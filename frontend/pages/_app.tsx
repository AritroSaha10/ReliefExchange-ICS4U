import Script from "next/script";
import "../styles/globals.css";
import { useEffect } from "react";
import { onAuthStateChanged, signOut } from "firebase/auth";
import auth from "@lib/firebase/auth";
import axios from "axios";
import convertBackendRouteToURL from "@lib/convertBackendRouteToURL";

export default function App({ Component, pageProps }) {
    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, user => {
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

        return () => unsubscribe()
    }, [])
    return (
        <>
            <Component {...pageProps} />
            <Script src="//translate.google.com/translate_a/element.js?cb=googleTranslateElementInit" defer />
        </>
    );
}