/**
 * @file File for the app component, which is used by Next.js to initialize pages.
 * More info about the component can be seen in its own dcumentation.
 * @author Aritro Saha
 */

import { useEffect } from "react";
import Script from "next/script";

import { onAuthStateChanged, signOut } from "firebase/auth";
import axios from "axios";

import auth from "@lib/firebase/auth";
import convertBackendRouteToURL from "@lib/convertBackendRouteToURL";

import "../styles/globals.css";

/**
 * App component, a part of Next.js that is used to initialize pages.
 * Attaches the Google Translate script, as well as kicking the user off the platform
 * if they get banned. 
 * 
 * The ban detection is in the App component so that they get notfied on the next 
 * route they navigate to once they get banned, instead of getting a random error
 * when trying to perform actions like looking at their profile.
 */
export default function App({ Component, pageProps }) {
    useEffect(() => {
        // Subscribe to authentication state changes
        const unsubscribe = onAuthStateChanged(auth, user => {
            // Only check ban state if user is logged in
            if (user && Object.keys(user).length !== 0) {
                // Kick the user off if they're banned
                const run = async () => {
                    const bannedRes = await axios.get(convertBackendRouteToURL(`/users/banned?uid=${user.uid}`))
                    if (bannedRes.data.banned) {
                        alert("You have been banned from our platform for breaking our rules. As such, you are not allowed to sign in.")
                        await signOut(auth)
                    }
                }
                run()
            }
        });

        // Unsubscribe function that runs once component is mounted
        return () => unsubscribe()
    }, [])
    return (
        <>
            <Component {...pageProps} />
            <Script src="//translate.google.com/translate_a/element.js?cb=googleTranslateElementInit" defer />
        </>
    );
}