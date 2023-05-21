
import { useState, useEffect } from "react";
import { useRouter } from "next/router";
import Image from "next/image";

import { User, onAuthStateChanged } from "firebase/auth";
import auth from "lib/firebase/auth";
import Layout from "@components/Layout";

import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime"
import axios from "axios";
import convertBackendRouteToURL from "lib/convertBackendRouteToURL";
import DonationCard from "@components/DonationCard";
import allTags from "lib/tag-types";
import Link from "next/link";
import Donation from "lib/types/donation";
dayjs.extend(relativeTime)

/**
 * Increase the resolution of a Google profile picture link. 
 * This is necessary since it's always 96px x 96px with no way to
 * request a higher resolution picture.
 * @param url The original image src link
 * @param newRes The new resolution in pixels (2-dimensional res not allowed)
 * @returns An altered src link with higher resolution
 */
const increasePFPResolution = (url: string, newRes: Number) => (
    url.replace("s96-c", `s${newRes.toFixed(0)}-c`)
)

export default function UserProfile() {
    const [loadingAuth, setLoadingAuth] = useState(true);
    const [user, setUser] = useState<User>(null);
    const [userData, setUserData] = useState<{ [key: string]: any }>();
    const [signedIn, setSignedIn] = useState<boolean>(false);

    const router = useRouter();

    /**
     * Refreshes user data on auth state change. Different from other versions of auth state change code,
     * as this gets the data of both the user and their donations.
     */
    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, newUser => {
            if (newUser && Object.keys(newUser).length !== 0) {
                // Set user data
                setUser(newUser);

                // Attempt to get the user's data from the backend
                axios.get(convertBackendRouteToURL(`/users/${newUser.uid}`)).then(async res => {
                    let data = res.data;
                    if (data.posts === null) {
                        data.posts = []
                    }

                    // Get data of each donation and replace the posts key with it
                    data.posts = await Promise.all(data.posts.map(async (post: { ID: string }) => {
                        // Get each donation
                        try {
                            const res = await axios.get(convertBackendRouteToURL(`/donations/${post.ID}`));

                            const donation: Donation = {
                                ...res.data,
                                creation_timestamp: new Date(res.data.creation_timestamp)
                            }

                            return donation;
                        } catch (e) {
                            console.log(e);
                            return null;
                        }
                    }));

                    data.posts = data.posts.filter((obj: any) => obj !== null)
                    data.posts.sort((a: Donation, b: Donation) => -(a.creation_timestamp.getTime() - b.creation_timestamp.getTime()))

                    setUserData(data);
                    setSignedIn(true);
                }).catch(err => {
                    console.error(err);
                    alert("Something went wrong while loading your user profile. Please try again later. Redirecting you to the home page...");
                    router.push("/");
                })
            } else {
                // User not logged in, throw them back to the front page
                setUser(null);
                setSignedIn(false);
                alert("You need to be signed in to access this page. Redirecting...");
                router.push("/");
            }

            setLoadingAuth(false);
        });

        return () => unsubscribe();
    }, [router]);

    if (!loadingAuth && signedIn) {
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
                                <p className="text-md text-white"><span className="font-semibold">Donations made: </span>{userData.donations_made}</p>
                            </div>
                        </div>

                        <div className="w-full md:w-2/3 lg:w-3/4">
                            <h2 className="text-3xl font-semibold text-white mb-2">Donations</h2>

                            {userData.posts.length !== 0 && (
                                <div className="flex flex-col self-center gap-4 lg:gap-6 w-full">
                                    {userData.posts.map(donation => {
                                        const tags = donation.tags ? donation.tags.map(tagName => allTags.find(tag => tag.name === tagName)) : []

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

                            {userData.posts.length === 0 && <span className="text-gray-200 text-md">It seems you haven&apos;t opened any donations yet. Click <Link href="/donations/create" className="text-blue-400 hover:underline active:text-blue-500">Donate</Link> in the navbar to make one!</span>}
                        </div>
                    </div>
                </div>
            </Layout>
        )
    } else {
        <span className="text-2xl text-center text-gray-200 font-semibold">Loading...</span>
    }
}