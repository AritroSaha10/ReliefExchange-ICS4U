/**
 * @file File for the profile page, which is accessible to logged-in users at /profile.
 * More info about the component can be seen in its own dcumentation.
 * @author Aritro Saha
 * @cite “Docs,” Docs | Next.js, https://nextjs.org/docs. 
 * @cite React, https://react.dev/. 
 * @cite “Get started with Firebase Authentication on websites,” Google, https://firebase.google.com/docs/auth/web/start. 
 * @cite “Day.js · 2kB javascript date utility library,” Day.js, https://day.js.org/en/. 
 */

import { useState, useEffect } from "react";
import { useRouter } from "next/router";
import Link from "next/link";
import Image from "next/image";

import { User, onAuthStateChanged, signOut } from "firebase/auth";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime"
import axios from "axios";

import Layout from "@components/Layout";
import DonationCard from "@components/DonationCard";

import convertBackendRouteToURL from "@lib/convertBackendRouteToURL";
import auth from "@lib/firebase/auth";
import allTags from "@lib/tag-types";
import Donation from "@lib/types/donation";
import Tag from "@lib/types/tag";
import increasePFPResolution from "@lib/increasePFPResolution";
import UserDataWithDonations from "@lib/types/userDataWithDonations"

dayjs.extend(relativeTime) // Allows us to calculate relative time easily

/**
 * The user profile page. Displays information about the user such as registration date and name,
 * as well as all the donations they've made. Only accessible whether they are signed in.
 */
export default function UserProfile() {
    // Whether we are currently waiting for user data from Firebase
    const [loadingAuth, setLoadingAuth] = useState(true);
    // Whether we are currently deleting the user's data
    const [deletingUser, setDeletingUser] = useState(false);
    // The authenticated user's data directly from Firebase auth 
    const [user, setUser] = useState<User>(null);
    // The user data from Firestore / the data we specically store
    const [userData, setUserData] = useState<UserDataWithDonations>();
    // Whether the user is signed in or not
    const [signedIn, setSignedIn] = useState<boolean>(false);

    // Router object to control current path state
    const router = useRouter();

    /**
     * Deletes all of the user's data from our website.
     */
    const deleteProfile = async () => {
        if (confirm("Are you sure you want to delete your account? You won't be able to recover it.")) {
            setDeletingUser(true)

            // Try sending a request to delete the user's data to the backend
            try {
                await axios.post(convertBackendRouteToURL(`/users/delete`), {
                    token: await user.getIdToken()
                })

                signOut(auth)

                alert("Your account has been deleted. Redirecting you to the home page...")
                router.push("/")
            } catch (e) {
                console.error(e)
                alert("Something went wrong while deleting your account. Please try again later.")
            }

            setDeletingUser(false)
        }
    }

    /**
     * Refreshes user data on auth state change. Different from other versions of auth state change code,
     * as this gets the data of both the user and their donations.
     */
    useEffect(() => {
        // Subscribes to authentication state change from Firebase
        const unsubscribe = onAuthStateChanged(auth, newUser => {
            // Only run if user is signed in
            if (newUser && Object.keys(newUser).length !== 0) {
                // Set user data
                setUser(newUser);

                // Attempt to get the user's data from the backend
                axios.get(convertBackendRouteToURL(`/users/${newUser.uid}`)).then(async res => {
                    let data = res.data;
                    // Empty arrays evaluate to null from Firebase, when
                    // they really should just be empty arrays.
                    if (data.posts === null) {
                        data.posts = []
                    }

                    // Get data of each donation and replace the post's key with it
                    data.posts = await Promise.all(data.posts.map(async (post: { ID: string }) => {
                        // Get each donation
                        try {
                            // Get raw donation
                            const res = await axios.get(convertBackendRouteToURL(`/donations/${post.ID}`));

                            // Convert the ISO string date to an actual date object 
                            const donation: Donation = {
                                ...res.data,
                                creation_timestamp: new Date(res.data.creation_timestamp)
                            }

                            return donation;
                        } catch (e) {
                            // Store a null object in the array, which we remove later
                            console.warn("Couldn't get a donation from database: ", e);
                            return null;
                        }
                    }));

                    // Filter out all the donations we couldn't fetch
                    data.posts = data.posts.filter((obj: any) => obj !== null)

                    // Sort by date descending
                    data.posts.sort((a: Donation, b: Donation) => -(a.creation_timestamp.getTime() - b.creation_timestamp.getTime()))

                    // Update state
                    setUserData(data);
                    setSignedIn(true);
                }).catch(err => {
                    // Let the user know of the error and push them back to their previous page
                    console.error(err);
                    alert("Something went wrong while loading your user profile. Please try again later. Redirecting you to the home page...");
                    router.back();
                })
            } else {
                // Let the user know and push them back to their previous page
                setUser(null);
                setSignedIn(false);
                alert("You need to be signed in to access this page. Redirecting...");
                router.back();
            }

            // Update state to reflect new changes in auth state
            setLoadingAuth(false);
        });

        // Unsubscribe from auth state changes on page unmount
        return () => unsubscribe();
    }, [router]);

    if (!loadingAuth && signedIn) {
        // Only show this page if the user is signed in
        return (
            <Layout name="Your Profile">
                <div className="p-10">
                    <h1 className="text-4xl font-semibold text-white mb-4 lg:mb-8">Your Account</h1>
                    <div className="flex flex-col md:flex-row w-full lg:p-8">
                        <div className="w-full md:w-1/3 lg:w-1/4">
                            <div className="flex flex-col">
                                <Image src={increasePFPResolution(user.photoURL, 400)} width={150} height={150} alt="user profile picture" className="rounded-full mb-2" />

                                <h3 className="text-2xl font-medium text-white">{user.displayName}</h3>
                                <p className="text-md text-white"><span className="font-semibold">Registered since:</span> {dayjs().to(user.metadata.creationTime)}</p>
                                <p className="text-md text-white"><span className="font-semibold">Last signed in: </span>{dayjs().to(user.metadata.lastSignInTime)}</p>
                                <p className="text-md text-white"><span className="font-semibold">Donations made: </span>{userData.donations_made.toString()}</p>
                                
                                <button
                                    className="flex items-center text-red-500 hover:text-red-600 active:text-red-700 disabled:text-red-900 duration-150"
                                    disabled={deletingUser}
                                    onClick={() => deleteProfile()}
                                >
                                    Delete Account
                                </button>
                            </div>
                        </div>

                        <div className="w-full md:w-2/3 lg:w-3/4">
                            <h2 className="text-3xl font-semibold text-white mb-2">Donations</h2>

                            {userData.posts.length !== 0 && (
                                <div className="flex flex-col self-center gap-4 lg:gap-6 w-full">
                                    {userData.posts.map(donation => {
                                        // Convert the tag names into tag objects
                                        const tags: Tag[] = (
                                            donation.tags ? donation.tags.map(tagName => allTags.find(tag => tag.name === tagName)) : []
                                        ).filter(tag => tag !== undefined);

                                        return (
                                            <DonationCard
                                                title={donation.title}
                                                date={donation.creation_timestamp}
                                                subtitle={donation.description}
                                                image={donation.img}
                                                tags={tags}
                                                href={`/donations/${donation.id}`}
                                                isAdmin={false} // Don't bother with showing admin data on their own posts
                                                reportCount={0}
                                                key={donation.id}
                                            />
                                        )
                                    })}
                                </div>
                            )}

                            {userData.posts.length === 0 && (
                                <span className="text-gray-200 text-md">
                                    It seems you haven&apos;t opened any donations yet. 
                                    Click 
                                    {" "}
                                    <Link href="/donations/create" className="text-blue-400 hover:underline active:text-blue-500">
                                        Donate
                                    </Link> 
                                    {" "}
                                    in the navbar to make one!
                                </span>
                            )}
                        </div>
                    </div>
                </div>
            </Layout>
        )
    } else {
        // Show some loading text as we wait for the user data to load
        return (
            <span className="text-2xl text-center text-gray-200 font-semibold my-4">Loading...</span>
        )
    }
}