import Image from "next/image";
import Link from "next/link";
import { useEffect, useState } from "react";
import { GetStaticPaths, GetStaticProps } from 'next'
import { useRouter } from "next/router";

import axios from "axios";
import { ParsedUrlQuery } from 'querystring'
import ReactMarkdown from "react-markdown";
import { User, onAuthStateChanged } from "firebase/auth";

import Layout from "@components/Layout";
import Donation from "@lib/types/donation";
import auth from "@lib/firebase/auth";
import UserData from "@lib/types/userData";
import RawDonation from "@lib/types/rawDonation";
import DonationWithUserData from "@lib/types/donationWithUserData";
import convertBackendRouteToURL from "@lib/convertBackendRouteToURL";

import { BiLeftArrowAlt } from "react-icons/bi"
import { FiFlag, FiTrash } from "react-icons/fi"
import { FaBan } from "react-icons/fa"


/**
 * The specific parameters that we need to get from the URL
 */
interface IParams extends ParsedUrlQuery {
    id: string
}

/**
 * Part of Next.js, gets all of the possible paths this page can have on the server.
 */
export const getStaticPaths: GetStaticPaths = async () => {
    // Get all the raw donations and extract the UIDs
    const rawDonations: RawDonation[] = (await axios.get(convertBackendRouteToURL("/donations/list"))).data
    const arr: string[] = rawDonations.map(donation => donation.id)

    return {
        paths: arr.map((id) => {
            return {
                params: { id },
            }
        }),
        fallback: "blocking" // If there's a new donation, run getStaticProps on it and cache its results instead of 404
    }
}

/**
 * Part of Next.js, fetches the data for a specific page on the server.
 */
export const getStaticProps: GetStaticProps = async (context) => {
    // Cast the parameters for the page request to our params interface
    const { id } = context.params as IParams

    try {
        // Get raw donation
        const rawDonation: RawDonation = (await axios.get(convertBackendRouteToURL(`/donations/${id}`))).data

        // Add user data to donation
        const rawDonationWithUserData = {
            ...rawDonation,
            owner: (await axios.get(convertBackendRouteToURL(`/users/${rawDonation.owner_id}`))).data
        }

        // Return it as a page prop, refreshing cache 1s after a page is served
        const props = { rawDonation: rawDonationWithUserData }
        return { props, revalidate: 1 }
    } catch (e) {
        if (e.response.status === 404) {
            // Donation doesn't exist
            return {
                notFound: true
            }
        } else {
            // Throw the error again for us to see in logs
            throw e
        }
    }
}

/**
 * The donation specific page, to show more information about a specific donation as well as allow them to contact the donator.
 */
export default function DonationSpecificPage({ rawDonation }) {
    const router = useRouter()

    // Convert the raw donation to an actual donation by converting the ISO string to
    // a date object. This is done since the Date object is not serializable and cannot
    // be sent in a JSON object as getStaticProps does.
    const donation: DonationWithUserData = {
        ...rawDonation,
        creation_timestamp: new Date(rawDonation.creation_timestamp)
    }

    // State vars
    const [user, setUser] = useState<User>(null);
    const [isAdmin, setIsAdmin] = useState(false);
    const [performingAction, setPerformingAction] = useState(false);

    /**
     * Sends a request to send a report regarding this post.
     */
    const sendReport = async () => {
        // Freeze other actions while performing this
        setPerformingAction(true)

        try {
            // Try sending a request to report
            await axios.post(convertBackendRouteToURL("/donations/report"), {
                donation_id: donation.id,
                token: await user.getIdToken()
            })

            alert("The post has been successfully reported. Thank you for helping us keep ReliefExchange clean.")
        } catch (e) {
            // Alert user of error and proceed
            if (e.response.status === 409) {
                alert("You have already reported this post. You cannot report it again.");
            } else {
                console.error(e)
                alert("Something went wrong while reporting this post. Please try again later.")
            }
        }

        // Unfreeze actions once done
        setPerformingAction(false)
    }

    /**
     * Sends a request to delete the post. Can only be run if the user is the owner of the post or an admin.
     */
    const deletePost = async () => {
        // Confirm with the user that they actually want to delete it
        if (confirm("Are you sure you want to delete this post? You won't be able to recover it.")) {
            // Freeze other actions
            setPerformingAction(true)
            
            try {
                // Try sending a request to delete, with proper authorization
                await axios.post(convertBackendRouteToURL(`/donations/${donation.id}/delete`), {}, {
                    headers: {
                        Authorization: `Bearer ${await user.getIdToken()}`,
                    },
                })

                // Redirect them back to front donations page
                alert("The post has been deleted. Redirecting you to the donations index page...")
                router.push("/donations")
            } catch (e) {
                // Alert user of issues
                console.error(e)
                alert("Something went wrong while deleting this post. Please try again later.")
            }

            // Unfreeze actions
            setPerformingAction(false)
        }
    }

    /**
     * Sends a request to ban the user.
     */
    const banUser = async () => {
        // Confirm with user to actually ban them
        if (confirm("Are you sure you want to ban this user? This will delete all their posts as well.")) {
            // Freeze donations
            setPerformingAction(true)

            try {
                // Try sending request to ban user
                alert("This will take a while. Please wait...")
                await axios.post(convertBackendRouteToURL(`/users/ban`), {
                    userToBan: donation.owner_id,
                    token: await user.getIdToken()
                })

                // Alert user of success and redirect back to donations home
                alert("The user has been banned. Redirecting you to the donations index page...")
                router.push("/donations")
            } catch (e) {
                // Alert user of error
                console.error(e)
                alert("Something went wrong while banning this user. Please try again later.")
            }

            // Unfreeze other actions
            setPerformingAction(false)
        }
    }

    /**
     * Add the event listener for auth state change
     */
    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, newUser => {
            // Only run if user is signed in
            if (newUser && Object.keys(newUser).length !== 0) {
                // Set user data
                setUser(newUser);

                // Check if user is admin by getting their user data
                try {
                    axios.get(convertBackendRouteToURL(`/users/${newUser.uid}`)).then(res => {
                        setIsAdmin(!!res.data.admin)
                    })
                } catch (e) {
                    // Silently record the error
                    console.error(e)
                }
            } else {
                // Not signed in, set to null
                setUser(null);
            }
        });

        // Unsubscribe from auth state changes on component dismount
        return () => unsubscribe();
    }, []);

    return (
        <Layout name={donation.title}>
            <div className="flex flex-col lg:items-center lg:justify-center flex-grow">
                <div className="flex flex-col lg:flex-row px-10 py-4 lg:px-20 lg:py-4 xl:px-60 xl:py-10 items-center">
                    <div className="flex flex-col gap-2 items-center lg:w-1/2">
                        <div className="flex gap-2 justify-between mb-2 w-full">
                            <Link
                                href="/donations"
                                className="flex items-center text-blue-500 hover:text-blue-600 active:text-blue-700 duration-150"
                            >
                                <BiLeftArrowAlt />
                                Back to Donations
                            </Link>

                            {user && user.uid !== donation.owner_id &&
                                <button
                                    className="flex items-center text-red-500 hover:text-red-600 active:text-red-700 disabled:text-red-900 duration-150"
                                    disabled={performingAction}
                                    onClick={() => sendReport()}
                                >
                                    Report
                                    <FiFlag className="ml-1" />
                                </button>
                            }

                            {user && (user.uid === donation.owner_id || isAdmin) &&
                                <button
                                    className="flex items-center text-red-500 hover:text-red-600 active:text-red-700 disabled:text-red-900 duration-150"
                                    disabled={performingAction}
                                    onClick={() => deletePost()}
                                >
                                    Delete
                                    <FiTrash className="ml-1" />
                                </button>
                            }
                        </div>

                        {donation.img ? <Image src={donation.img} alt="Featured image" height={500} width={500} className="rounded-md object-cover object-center" /> : <></>}

                        {user && isAdmin && (
                            <div className="flex gap-2 justify-between mb-2 w-full">
                                <button
                                    className="flex items-center text-red-500 hover:text-red-600 active:text-red-700 disabled:text-red-900 duration-150"
                                    disabled={performingAction}
                                    onClick={() => banUser()}
                                >
                                    <FaBan className="mr-1" />
                                    Ban User
                                </button>

                                <span
                                    className="flex items-center text-orange-500"
                                >
                                    Reports: {donation.reports.length}
                                    <FiFlag className="ml-1" />
                                </span>
                            </div>
                        )}
                    </div>

                    <div className="lg:ml-5 flex flex-col items-center lg:items-start">
                        <h1 className="text-white text-4xl font-semibold text-center break-all">{donation.title}</h1>
                        <h3 className="text-gray-300 text-md text-center break-all">Posted on {donation.creation_timestamp.toLocaleDateString("en-US", { day: "numeric", month: "long", year: "numeric", })}</h3>
                        <h3 className="text-gray-300 text-md text-center break-all">Available in &quot;{donation.location}&quot;</h3>

                        <br />

                        {/* This makes sure that all external links in the description open another tab */}
                        <base target="_blank" />

                        <ReactMarkdown className="text-white break-all">{donation.description}</ReactMarkdown>

                        <br />

                        <a className="py-2 px-4 bg-blue-500 font-semibold text-center text-white rounded-lg hover:bg-blue-600 duration-75" href={`mailto:${donation.owner.email}`}>Contact {donation.owner.display_name} for More Info</a>
                    </div>
                </div>
            </div>
        </Layout>
    )
}
