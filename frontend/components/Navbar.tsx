import Link from "next/link";
import { useState, useEffect } from "react";
import { onAuthStateChanged } from "firebase/auth";
import auth from "lib/firebase/auth";

import { GoThreeBars } from "react-icons/go"

const baseLinks = [
    {
        name: "Home",
        link: "/",
        id: "home",
        priority: false
    },
    {
        name: "Donations",
        link: "/donations",
        id: "donations",
        priority: false
    },
];

const signedInLinks = [
    {
        name: "Donate",
        link: "/donations/create",
        id: "donate",
        priority: false
    },
    {
        name: "Profile",
        link: "/profile",
        id: "profile",
        priority: false
    },
    {
        name: "Sign Out",
        link: "/sign-out",
        id: "sign-out",
        priority: true
    },
];

const signedOutLinks = [
    {
        name: "Sign In / Register",
        link: "/sign-in",
        id: "sign-in",
        priority: true
    },
];

export default function Header() {
    const [showDropdown, setShowDropdown] = useState(false);
    const [isSignedIn, setSignedIn] = useState(false);

    // Refresh user data state on auth state change
    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, user => {
            setSignedIn(user !== null);
        });

        return () => unsubscribe();
    }, []);

    return (
        <header className="bg-slate-800 py-2 lg:py-4 sticky">
            <div className="container px-4 mx-auto lg:flex lg:items-center">
                <div className="flex justify-between items-center">
                    <Link href="/">
                        {/* <Image src={Logo} alt="logo" width={50} height={50} /> */}
                        <span className="text-2xl font-mono tracking-wider font-bold text-white">ReliefExchange</span>
                    </Link>

                    <button
                        className="border border-solid border-gray-600 px-3 py-1 rounded text-gray-600 opacity-50 hover:opacity-75 lg:hidden"
                        aria-label="Menu"
                        data-test-id="navbar-menu"
                        onClick={
                            () => {
                                setShowDropdown(!showDropdown);
                            }}
                    >
                        <GoThreeBars />
                    </button>
                </div>

                <div className={`${showDropdown ? "flex" : "hidden"} lg:flex flex-col lg:flex-row lg:ml-auto mt-3 lg:mt-0`} data-test-id="navbar">
                    {
                        baseLinks.map(({ name, link, priority, id }) =>
                            <Link key={name} href={link} className={`${priority ? "text-blue-600 hover:bg-blue-600 hover:text-white text-center border border-solid border-blue-600 mt-1 lg:mt-0 lg:ml-1" : "text-gray-300 hover:bg-gray-200 hover:text-gray-700 "} p-2 lg:px-4 lg:mx-2 rounded duration-300 transition-colors `}
                                data-test-id={`navbar-${id}`}>
                                {name}
                            </Link>
                        )
                    }

                    {
                        isSignedIn ? signedInLinks.map(({ name, link, priority, id }) =>
                            <Link key={name} href={link} className={`${priority ? "text-blue-600 hover:bg-blue-600 hover:text-white text-center border border-solid border-blue-600 mt-1 lg:mt-0 lg:ml-1" : "text-gray-300 hover:bg-gray-200 hover:text-gray-700 "} p-2 lg:px-4 lg:mx-2 rounded duration-300 transition-colors `}
                                data-test-id={`navbar-${id}`}>
                                {name}
                            </Link>
                        ) : signedOutLinks.map(({ name, link, priority, id }) =>
                            <Link key={name} href={link} className={`${priority ? "text-blue-600 hover:bg-blue-600 hover:text-white text-center border border-solid border-blue-600 mt-1 lg:mt-0 lg:ml-1" : "text-gray-300 hover:bg-gray-200 hover:text-gray-700 "} p-2 lg:px-4 lg:mx-2 rounded duration-300 transition-colors `}
                                data-test-id={`navbar-${id}`}>
                                {name}
                            </Link>
                        )
                    }
                </div>
            </div>
        </header>
    )
}