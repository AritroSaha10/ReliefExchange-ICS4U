import Image from "next/image";
import { useEffect, useState } from "react"

import Layout from "@components/Layout";
import auth from "lib/firebase/auth";

import { onAuthStateChanged, signOut } from "firebase/auth";

import { useRouter } from "next/router";

export default function SignOut() {
    const router = useRouter();
    
    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, user => {
            if (!user) {
                alert("You are signed out. Redirecting...");
                router.push("/");
            } else {
                signOut(auth).catch((e) => {
                    alert("Something went wrong while signing you out. Please try again.");
                    console.error(e);
                });
            }
        });

        return () => unsubscribe();
    }, []);
    
    return (
        <Layout name="Sign In">
            <div className="flex flex-col flex-grow items-center justify-center">
                <div className="flex flex-col items-center justify-center gap-2">
                    <h1 className="text-3xl font-semibold text-white mb-4">Sign out</h1>
                    <h1 className="text-xl text-white mb-4">Logging you out...</h1>
                </div>
            </div>
        </Layout>
    )
}