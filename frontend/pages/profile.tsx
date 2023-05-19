
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
dayjs.extend(relativeTime)

export default function UserProfile() {
    const [loadingAuth, setLoadingAuth] = useState(true);
    const [user, setUser] = useState<User>(null);
    const [userData, setUserData] = useState<{ [key: string]: any }>();
    const [signedIn, setSignedIn] = useState<boolean>(false);

    const router = useRouter();

    const increasePFPResolution = (url: string, newRes: Number) => (
        url.replace("s96-c", `s${newRes.toFixed(0)}-c`)
    )

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
                            return res.data;
                        } catch (e) {
                            console.log(e);
                            return null;
                        }
                    }));

                    data.posts = data.posts.filter((obj: any) => obj !== null)

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
    }, []);

    if (!loadingAuth && signedIn) {
        return (
            <Layout name="Your Profile">
                <div className="p-10">
                    <h1 className="text-4xl font-semibold text-white mb-4">Your Account</h1>
                    <div className="flex flex-col md:flex-row w-full">
                        <div className="w-full md:w-1/3 lg:w-1/4 bg-gray-800">
                            <div className="flex flex-col items-center text-center">
                                <Image src={increasePFPResolution(user.photoURL, 400)} width={150} height={150} alt="user profile picture" className="rounded-full" />
                                <h3 className="text-2xl font-medium text-white">{user.displayName}</h3>
                                <p className="text-md text-white"><span className="font-semibold">Registered since:</span> {dayjs().to(user.metadata.creationTime)}</p>
                                <p className="text-md text-white"><span className="font-semibold">Last signed in: </span>{dayjs().to(user.metadata.lastSignInTime)}</p>
                                <p className="text-md text-white"><span className="font-semibold">Donations made: </span>{userData.donations_made}</p>
                            </div>
                        </div>

                        <div className="w-full md:w-2/3 lg:w-3/4 bg-slate-700">
                            <h2 className="text-3xl font-semibold text-white">Donations</h2>
                        </div>
                    </div>
                </div>
            </Layout>
        )
    } else {
        <span className="text-2xl text-center text-gray-200 font-semibold">Loading...</span>
    }
}