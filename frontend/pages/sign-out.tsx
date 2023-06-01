/**
 * @file File for the sign-out page, which is accessible to logged-in users at /sign-out.
 * @author Aritro Saha
 */

import { useEffect } from "react"
import { useRouter } from "next/router";

import { signOut } from "firebase/auth";

import Layout from "@components/Layout";
import auth from "@lib/firebase/auth";

/**
 * The sign-out page. Signs the user out.
 */
export default function SignOut() {
    // Necessary hook for redirecting user
    const router = useRouter();

    /**
     * Sign them out once we know they're signed in.
     */
    useEffect(() => {
        signOut(auth).catch((e) => {
            alert("Something went wrong while signing you out. Please try again.");
            console.error(e);
        }).finally(() => {
            router.push("/")
        });
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

    return (
        <Layout name="Sign Out">
            <div className="flex flex-col flex-grow items-center justify-center">
                <div className="flex flex-col items-center justify-center gap-2">
                    <h1 className="text-3xl font-semibold text-white mb-4">Sign out</h1>
                    <h1 className="text-xl text-white mb-4">Logging you out...</h1>
                </div>
            </div>
        </Layout>
    )
}