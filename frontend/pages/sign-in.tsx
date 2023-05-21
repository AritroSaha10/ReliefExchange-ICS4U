import Image from "next/image";
import { useEffect } from "react"

import Layout from "@components/Layout";
import auth from "lib/firebase/auth";

import { GoogleAuthProvider, setPersistence, signInWithPopup, browserLocalPersistence, onAuthStateChanged } from "firebase/auth";

import GoogleLogo from "@media/social-media-logos/google.png";
import MicrosoftLogo from "@media/social-media-logos/microsoft.png";
import FacebookLogo from "@media/social-media-logos/facebook.png";
import { useRouter } from "next/router";
import axios from "axios";
import convertBackendRouteToURL from "lib/convertBackendRouteToURL";

export default function SignIn() {
    const provider = new GoogleAuthProvider();
    const router = useRouter();

    /**
     * Prompt the user to sign in 
     */
    const continueWithGoogle = async () => {
        try {
            // Set persistence
            await setPersistence(auth, browserLocalPersistence);

            const res = await signInWithPopup(auth, provider);

            if (res.user.metadata.creationTime == res.user.metadata.lastSignInTime) {
                try {
                    await axios.post(convertBackendRouteToURL("/users/new"), {
                        token: await res.user.getIdToken()
                    })
                } catch (e) {
                    console.error(e)
                }
            }

            console.log(res);
        } catch (e) {
            console.error(e);
            alert("Something went wrong. Please try again.");
        }
    }

    /**
     * Bring them back to front page once they're signed in.
     */
    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, user => {
            if (user) {
                alert("You are signed in! Redirecting...");
                router.push("/");
            }
        });

        return () => unsubscribe();
    }, [router]);

    return (
        <Layout name="Sign In">
            <div className="flex flex-col flex-grow items-center justify-center">
                <div className="flex flex-col items-center justify-center gap-2">
                    <h1 className="text-3xl font-semibold text-white mb-4">Log In / Sign Up</h1>
                    <button className="flex flex-row gap-3 justify-center items-end bg-slate-700 hover:bg-slate-600 active:bg-slate-800 duration-150 px-5 py-3 rounded-lg w-full" onClick={() => { continueWithGoogle() }}>
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