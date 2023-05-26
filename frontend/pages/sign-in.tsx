import Image from "next/image";
import { useEffect, useState } from "react"

import Layout from "@components/Layout";
import auth from "@lib/firebase/auth";

import { GoogleAuthProvider, setPersistence, browserLocalPersistence, onAuthStateChanged, getRedirectResult, signInWithRedirect } from "firebase/auth";

import GoogleLogo from "@media/social-media-logos/google.png";
import { useRouter } from "next/router";
import axios from "axios";
import convertBackendRouteToURL from "lib/convertBackendRouteToURL";

/**
 * The sign-in page. Handles the sign-in flow. Only accessible if the user isn't signed in.
 */
export default function SignIn() {
    // State and hooks necessary for sign-in.
    const router = useRouter();
    const [signingIn, setSigningIn] = useState(false)

    /**
     * Prompt the user to sign in using their Google account.
     */
    const continueWithGoogle = async () => {
        try {
            // Toggle sign-in so they can't click the button multiple times
            setSigningIn(true);

            // Set the auth persistance so they don't have to sign in multiple times
            await setPersistence(auth, browserLocalPersistence);

            // Set up authentication provider (in this case, Google)
            const provider = new GoogleAuthProvider();
            provider.addScope("email");

            // Redirect the user to sign-in page
            await signInWithRedirect(auth, provider);
        } catch (e) {
            // Log the error and let the user know
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
            // Create the provider class just like it was created when starting the log in flow
            const provider = new GoogleAuthProvider();
            provider.addScope("email");

            // Manage the results if there was a redirect
            const res = await getRedirectResult(auth);
            if (res !== null) {
                // Disable sign-in button to prevent them starting another login flow
                setSigningIn(true);

                // Actually a redirect, handle sign-in
                if (res.user.metadata.creationTime === res.user.metadata.lastSignInTime) {
                    try {
                        // Create a document for their user data in our Firestore DB if user is new
                        await axios.post(convertBackendRouteToURL("/users/new"), {
                            token: await res.user.getIdToken()
                        })
                    } catch (e) {
                        console.error(e)
                    }
                }

                // Redirect them to the home page
                router.push("/");
            } else {
                // No previous log-in flow, let them use the page as normal (if they aren't logged in)

                // While we should technically be providing an unsubscribe function to React,
                // it's not possible since this is in an async function. It must be in this async
                // function since we first must wait for the result from the redirect
                // before subscribing, or else we'll get race conditions! It's fine if we don't unsubscribe here
                // since this will only run once anyways. I think...
                onAuthStateChanged(auth, user => {
                    if (user) {
                        // Redirect to home page if they're already signed in
                        alert("You are already signed in! Redirecting...");
                        router.push("/");
                    }
                });
            }
        })();
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

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