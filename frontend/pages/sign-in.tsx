import Image from "next/image";
import { useEffect, useState } from "react"

import Layout from "@components/Layout";
import auth from "@lib/firebase/auth";

import { GoogleAuthProvider, setPersistence, signInWithPopup, browserLocalPersistence, onAuthStateChanged, getRedirectResult, signInWithRedirect } from "firebase/auth";

import GoogleLogo from "@media/social-media-logos/google.png";
import MicrosoftLogo from "@media/social-media-logos/microsoft.png";
import FacebookLogo from "@media/social-media-logos/facebook.png";
import { useRouter } from "next/router";
import axios from "axios";
import convertBackendRouteToURL from "lib/convertBackendRouteToURL";

export default function SignIn() {
    const router = useRouter();
    const [signingIn, setSigningIn] = useState(false)

    /**
     * Prompt the user to sign in 
     */
    const continueWithGoogle = async () => {
        try {
            setSigningIn(true);
            // Set persistence
            await setPersistence(auth, browserLocalPersistence);

            // Redirect the user to sign-in
            const provider = new GoogleAuthProvider();
            provider.addScope("email");
            await signInWithRedirect(auth, provider);
        } catch (e) {
            console.error(e);
            alert("Something went wrong. Please try again.");
        }
    }

    /**
     * Bring them back to front page once they're signed in.
     */
    useEffect(() => {
        // Make sure to run this when the user first signs in
        (async () => {
            const provider = new GoogleAuthProvider();
            provider.addScope("email");

            const res = await getRedirectResult(auth);
            if (res !== null) {
                setSigningIn(true);

                // Actually a redirect, handle sign-in
                if (res.user.metadata.creationTime == res.user.metadata.lastSignInTime) {
                    try {
                        await axios.post(convertBackendRouteToURL("/users/new"), {
                            token: await res.user.getIdToken()
                        })
                    } catch (e) {
                        console.error(e)
                    }
                }

                router.push("/");
            } else {
                // While we should technically be providing an unsubscribe function to React,
                // it's not possible since this is in an async function. It must be in this async
                // function since we first must wait for the result from the redirect
                // before subscribing, or else we'll get race conditions! It's fine if we don't unsubscribe here
                // since this will only run once anyways. I think...
                onAuthStateChanged(auth, user => {
                    if (user) {
                        alert("You are already signed in! Redirecting...");
                        router.push("/");
                    }
                });
            }
        })();
    }, []);

    return (
        <Layout name="Sign In">
            <div className="flex flex-col flex-grow items-center justify-center">
                <div className="flex flex-col items-center justify-center gap-2">
                    <h1 className="text-3xl font-semibold text-white mb-4">Log In / Sign Up</h1>
                    <button className="flex flex-row gap-3 justify-center items-end bg-slate-700 hover:bg-slate-600 active:bg-slate-800 duration-150 px-5 py-3 rounded-lg w-full disabled:text-gray-200 disabled:bg-gray-800" onClick={() => { continueWithGoogle() }} disabled={signingIn}>
                        <Image src={GoogleLogo} width={25} height={25} alt="" />
                        <span className="text-xl text-white">
                            Continue with Google
                        </span>
                    </button>
                </div>
            </div>
        </Layout>
    )
}