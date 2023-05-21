import axios from "axios";
import Image from "next/image";

import { GetStaticPaths, GetStaticProps } from 'next'
import { ParsedUrlQuery } from 'querystring'
import Donation from "lib/types/donation";
import Layout from "@components/Layout";
import { ReactMarkdown } from "react-markdown/lib/react-markdown";

import UserData from "lib/types/userData";
import RawDonation from "lib/types/rawDonation";
import DonationWithUserData from "lib/types/donationWithUserData";

import { BiLeftArrowAlt } from "react-icons/bi"
import { FiFlag, FiTrash } from "react-icons/fi"
import { FaBan } from "react-icons/fa"

import Link from "next/link";
import { useEffect, useState } from "react";
import { User, onAuthStateChanged } from "firebase/auth";
import { useRouter } from "next/router";
import auth from "@lib/firebase/auth";
import convertBackendRouteToURL from "lib/convertBackendRouteToURL";

interface IParams extends ParsedUrlQuery {
    id: string
}

export const getStaticPaths: GetStaticPaths = async () => {
    const rawDonations: RawDonation[] = (await axios.get(convertBackendRouteToURL("/donations/list"))).data
    const arr: string[] = rawDonations.map(donation => donation.id)

    return {
        paths: arr.map((id) => {
            return {
                params: { id },
            }
        }),
        fallback: "blocking"
    }
}

export const getStaticProps: GetStaticProps = async (context) => {
    const { id } = context.params as IParams

    try {
        const rawDonation: RawDonation = (await axios.get(convertBackendRouteToURL(`/donations/${id}`))).data
        const rawDonationWithUserData = {
            ...rawDonation,
            owner: (await axios.get(convertBackendRouteToURL(`/users/${rawDonation.owner_id}`))).data
        }

        const props = { rawDonation: rawDonationWithUserData }
        return { props, revalidate: 1 }
    } catch (e) {
        if (e.response.status === 404) {
            return {
                notFound: true
            }
        } else {
            throw e
        }
    }
}

export default function DonationSpecificPage({ rawDonation }) {
    const router = useRouter()
    const donation: DonationWithUserData = {
        ...rawDonation,
        creation_timestamp: new Date(rawDonation.creation_timestamp)
    }

    const [user, setUser] = useState<User>(null);
    const [isAdmin, setIsAdmin] = useState(false);
    const [performingAction, setPerformingAction] = useState(false);

    /**
     * Sends a request to send a report regarding this post.
     */
    const sendReport = async () => {
        setPerformingAction(true)

        try {
            await axios.post(convertBackendRouteToURL("/donations/report"), {
                donation_id: donation.id,
                token: await user.getIdToken()
            })

            alert("The post has been successfully reported. Thank you for helping us keep ReliefExchange clean.")
        } catch (e) {
            if (e.response.status === 409) {
                alert("You have already reported this post. You cannot report it again.");
            } else {
                console.error(e)
                alert("Something went wrong while reporting this post. Please try again later.")
            }
        }

        setPerformingAction(false)
    }

    /**
     * Sends a request to delete the post. Can only be run if the user is the owner of the post or an admin.
     */
    const deletePost = async () => {
        if (confirm("Are you sure you want to delete this post? You won't be able to recover it.")) {
            setPerformingAction(true)

            try {
                await axios.post(convertBackendRouteToURL(`/donations/${donation.id}/delete`), {}, {
                    headers: {
                        Authorization: `Bearer ${await user.getIdToken()}`,
                    },
                })

                alert("The post has been deleted. Redirecting you to the donations index page...")
                router.push("/donations")
            } catch (e) {
                console.error(e)
                alert("Something went wrong while deleting this post. Please try again later.")
            }

            setPerformingAction(false)
        }
    }

    /**
     * Sends a request to ban the user.
     */
    const banUser = async () => {
        if (confirm("Are you sure you want to ban this user? This will delete all their posts as well.")) {
            setPerformingAction(true)

            try {
                alert("This will take a while. Please wait...")
                await axios.post(convertBackendRouteToURL(`/users/ban`), {
                    userToBan: donation.owner_id,
                    token: await user.getIdToken()
                })

                alert("The user has been banned. Redirecting you to the donations index page...")
                router.push("/donations")
            } catch (e) {
                console.error(e)
                alert("Something went wrong while banning this user. Please try again later.")
            }

            setPerformingAction(false)
        }
    }

    /**
     * Add the event listener for auth state change
     */
    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, newUser => {
            if (newUser && Object.keys(newUser).length !== 0) {
                // Set user data
                setUser(newUser);

                try {
                    axios.get(convertBackendRouteToURL(`/users/${newUser.uid}`)).then(res => {
                        setIsAdmin(!!res.data.admin)
                    })
                } catch (e) {
                    console.error(e)
                }
            } else {
                setUser(null);
            }
        });

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

                        <ReactMarkdown skipHtml={true} className="text-white break-all">{donation.description}</ReactMarkdown>

                        <br />

                        <a className="py-2 px-4 bg-blue-500 font-semibold text-center text-white rounded-lg hover:bg-blue-600 duration-75" href={`mailto:${donation.owner.email}`}>Contact {donation.owner.display_name} for More Info</a>
                    </div>
                </div>
            </div>
        </Layout>
    )
}
