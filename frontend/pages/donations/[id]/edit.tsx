/**
 * @file File for the donation edit page, which is only accessible to either admins
 * or the author of the donation.
 * @author Aritro Saha
 * @cite React, https://react.dev/. 
 * @cite D. Omotayo, “How to implement ReCAPTCHA in a React application,” LogRocket Blog, https://blog.logrocket.com/implement-recaptcha-react-application/. 
 * @cite “Docs,” Docs | Next.js, https://nextjs.org/docs. 
 */

import { useState, useEffect, useRef, FormEventHandler } from "react";
import { GetStaticPaths, GetStaticProps } from "next";
import { useRouter } from "next/router";
import dynamic from "next/dynamic";

import axios from "axios";
import { getIdToken, onAuthStateChanged, User } from "firebase/auth";
import ReCAPTCHA from "react-google-recaptcha"
import Multiselect from 'multiselect-react-dropdown';
import { ParsedUrlQuery } from "querystring";
import * as commands from "@uiw/react-md-editor/lib/commands";

import Layout from "@components/Layout";
import auth from "@lib/firebase/auth";
import allTags from "@lib/tag-types";
import RawDonation from "@lib/types/rawDonation";
import convertBackendRouteToURL from "@lib/convertBackendRouteToURL";

import "@uiw/react-md-editor/markdown-editor.css";
import "@uiw/react-markdown-preview/markdown.css";

// Don't try to render this component on the server
const MDEditor = dynamic(
    () => import("@uiw/react-md-editor").then((mod) => mod.default),
    { ssr: false }
);

/**
 * Donation Edit page, where signed-in users can edit donations that show up on the index.
 */
export default function EditDonation({ originalDonation }: { originalDonation: RawDonation }) {
    // Necessary hooks
    const router = useRouter();
    const captchaRef = useRef(null);

    // All state variables for form elements and user data
    const [loadingAuth, setLoadingAuth] = useState(true);
    const [user, setUser] = useState<User>(null);
    const [signedIn, setSignedIn] = useState<boolean>(false);
    const [tagsSelected, setTagsSelected] = useState(
        originalDonation.tags.map(key => allTags.find(
            obj => obj.name === key
        ))
    );
    const [descriptionMD, setDescriptionMD] = useState(originalDonation.description);
    const [submitting, setSubmitting] = useState(false);

    // Site key for using Google ReCAPTCHA v2
    const RECAPTCHA_SITE_KEY = process.env.NEXT_PUBLIC_RECAPTCHA_SITE_KEY;

    /**
     * Refresh user data on auth change
     */
    useEffect(() => {
        // Subscribe to auth state changes
        const unsubscribe = onAuthStateChanged(auth, newUser => {
            // Only run if user is signed-in
            if (newUser && Object.keys(newUser).length !== 0) {
                // Check whether user is admin to ensure they are allowed to access this page
                axios.get(convertBackendRouteToURL(`/users/admin?uid=${newUser.uid}`)).then(res => {
                    // Don't allow if they're not an admin or original author
                    if (newUser.uid !== originalDonation.owner_id && !res.data["admin"]) {
                        alert("You cannot edit this post, as you are not its author. Redirecting...")
                        router.push("/")
                    } else {
                        // Set user data
                        setUser(newUser);
                        setSignedIn(true);
                    }
                }).catch(err => {
                    console.error(err)
                    alert("Something went wrong. Please try again.")
                    router.push("/")
                })

            } else {
                setUser(null);
                setSignedIn(false);

                // Throw user back to homepage
                alert("You need to be signed in to access this page. Redirecting...");
                router.push("/");
            }

            // Update state
            setLoadingAuth(false);
        });

        // Unsubscribe from auth state changes on page dismount to avoid zombie auth state subscribers
        return () => unsubscribe();
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

    /**
     * Handle the submit event for the main form
     * @param e Event handler data for a form
     * @returns Nothing.
     */
    const onSubmit: FormEventHandler<HTMLFormElement> = async (e) => {
        // Convert form data to keys and don't refresh page
        const formData = Object.fromEntries((new FormData(e.currentTarget)).entries());
        e.preventDefault();

        // Freeze submit button
        setSubmitting(true);

        // Confirm the CAPTCHA with the server
        const token = captchaRef.current.getValue();
        captchaRef.current.reset();
        if (!token) {
            alert("Please complete the CAPTCHA and try again.");
            setSubmitting(false);
            return;
        }
        try {
            // Send request to server to get result of CAPTCHA
            const res = await axios.post(convertBackendRouteToURL(`/confirmCAPTCHA?token=${token}`))
            if (!res.data.human) throw "User was detected to be a bot by ReCAPTCHA."
        } catch (e) {
            // Handle error, let user know of problem
            alert("Something went wrong. Please try again.");
            console.error(e);
            setSubmitting(false);

            // Don't proceed with rest of donation adding process
            return
        }

        // CAPTCHA confirmed, now proceed with rest
        // Convert the current timestamp to UTC for a consistent timezone
        const date = new Date();
        const nowUTC = new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(),
            date.getUTCDate(), date.getUTCHours(),
            date.getUTCMinutes(), date.getUTCSeconds()));

        // Image uploaded, now prepare the data to send to endpoint
        const idToken = await getIdToken(user, true);
        const donationData = {
            "title": formData["product-name"],
            "description": descriptionMD,
            "location": formData["product-location"],
            "img": "",
            "tags": tagsSelected.map(obj => obj.name),
            "creation_timestamp": nowUTC.toISOString(),
            "ownerID": user.uid
        };

        // Send the prep'd data to our endpoint
        try {
            const apiRes = await axios.post(convertBackendRouteToURL("/donations/edit"), {
                data: donationData,
                token: idToken,
                id: originalDonation.id
            });
            alert("Your donation was successfully edited! Redirecting you to its page...");
            router.push(`/donations/${apiRes.data}`);
        } catch (e) {
            // Let user know of issue
            alert("Something went wrong while submitting your donation. Please try again.");
            console.error(e);
        }

        // Unfreeze submit button
        setSubmitting(false);
    }

    return (
        <Layout name="Edit a Donation">
            <div className="flex flex-col gap-4 p-10 flex-grow min-w-screen">
                <div className="flex flex-col gap-2 mb-4">
                    <h1 className="text-4xl text-white font-bold text-center">Edit Donation Offer</h1>
                    <p className="text-md text-gray-200 text-center">Edit a donation offer.</p>
                </div>
                {loadingAuth && (
                    <p className="text-center text-gray-200 text-md">Loading...</p>
                )}
                {signedIn && (
                    <form onSubmit={onSubmit}>
                        <div className="flex flex-col gap-6">
                            <div className="flex flex-col items-center">
                                <h3 className="text-white text-2xl font-medium mb-2 text-center lg:text-left">Product Name: (Max. 100 characters) <span className="text-red-500"> *</span></h3>
                                <div className="flex flex-col md:flex-row gap-4">
                                    <input
                                        name="product-name"
                                        type="text"
                                        placeholder={`Product Name...`}
                                        className="rounded py-2 px-3 w-60 sm:w-72 align-middle text-gray-700 outline-none ring-2 ring-blue-100 focus:ring-blue-300 duration-200"
                                        required
                                        defaultValue={originalDonation.title}
                                        maxLength={100}
                                    />
                                </div>
                            </div>

                            <div className="flex flex-col items-center">
                                <h3 className="text-white text-2xl font-medium mb-2 text-center lg:text-left">Product Description: (Max. 1000 characters) <span className="text-red-500"> *</span></h3>
                                <div className="flex flex-col gap-4 w-full items-center" data-color-mode="light">
                                    <MDEditor
                                        value={descriptionMD}
                                        onChange={(val) => {
                                            // Enforce character limit
                                            if (val.length > 1000) {
                                                val = val.slice(0, 1000)
                                            }

                                            setDescriptionMD(val)
                                        }}
                                        className="w-full lg:w-2/3"
                                        commands={[
                                            commands.bold, commands.hr, commands.italic, commands.divider, commands.codeEdit, commands.codeLive, commands.codePreview, commands.divider,
                                            commands.fullscreen,
                                        ]}
                                    />
                                </div>
                            </div>

                            <div className="flex flex-col items-center">
                                <h3 className="text-white text-2xl font-medium mb-2 text-center lg:text-left">Product Tags: (Max. 3 tags) <span className="text-red-500"> *</span></h3>
                                <div className="flex flex-col gap-4 w-2/3 lg:w-1/2">
                                    <Multiselect
                                        options={allTags} // Options to display in the dropdown
                                        selectedValues={tagsSelected} // Preselected value to persist in dropdown
                                        onSelect={setTagsSelected} // Function will trigger on select event
                                        onRemove={setTagsSelected} // Function will trigger on remove event
                                        displayValue="name" // Property name to display in the dropdown options
                                        selectionLimit={3}
                                        className="bg-white rounded-xl border-none transition-all duration-150"
                                    />
                                </div>
                            </div>

                            <div className="flex flex-col items-center">
                                <h3 className="text-white text-2xl font-medium mb-2 text-center lg:text-left">Location: (Max. 100 characters) <span className="text-red-500"> *</span></h3>
                                <div className="flex flex-col md:flex-row gap-4">
                                    <input
                                        type="text"
                                        name="product-location"
                                        placeholder={`Location...`}
                                        className="rounded py-2 px-3 w-72 sm:w-80 align-middle text-gray-700 outline-none ring-2 ring-blue-100 focus:ring-blue-300 duration-200"
                                        required
                                        defaultValue={originalDonation.location}
                                        maxLength={100}
                                    />
                                </div>
                            </div>

                            <div className="flex flex-col items-center gap-2 self-center mb-4">
                                <ReCAPTCHA sitekey={RECAPTCHA_SITE_KEY} ref={captchaRef} />

                                <label>
                                    <button
                                        className={`flex items-center bg-green-500 hover:bg-green-700 text-white font-semibold py-4 px-8 rounded text-xl capitalize duration-75 ${submitting && "bg-green-800 cursor-default"}`}
                                        disabled={submitting}
                                    >
                                        <svg className={`animate-spin ml-1 mr-3 h-5 w-5 text-white ${submitting ? "inline-block" : "hidden"}`} xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                        </svg>

                                        Submit
                                    </button>
                                    <input type="submit" value="Submit" hidden disabled={submitting} />
                                </label>

                                <p className="text-sm text-gray-300 text-center w-3/4 md:w-1/2">
                                    ReliefExchange is not responsible for products exchanged on our platform. By submitting the form, you acknowledge and agree to this.
                                </p>
                            </div>
                        </div>
                    </form>
                )}
            </div>
        </Layout>
    )
}

/**
 * Part of Next.js, gets all of the possible paths this page can have on the server.
 */
export const getStaticPaths: GetStaticPaths = async () => {
    // Get all the raw donations and extract the UIDs
    const rawDonations: RawDonation[] = (await axios.get(convertBackendRouteToURL("/donations/list"))).data
    const arr: string[] = rawDonations.map(donation => donation.id)

    return {
        paths: arr.map((id) => {
            return {
                params: { id },
            }
        }),
        fallback: "blocking" // If there's a new donation, run getStaticProps on it and cache its results instead of 404
    }
}

/**
 * Part of Next.js, fetches the data for a specific page on the server.
 */
export const getStaticProps: GetStaticProps = async (context) => {
    // Cast the parameters for the page request to our params interface
    const { id } = context.params as IParams

    try {
        // Get raw donation
        const rawDonation: RawDonation = (await axios.get(convertBackendRouteToURL(`/donations/${id}`))).data

        // Return it as a page prop, refreshing cache 1s after a page is served
        const props = { originalDonation: rawDonation }
        console.log(props)
        return { props, revalidate: 1 }
    } catch (e) {
        if (e.response.status === 404) {
            // Donation doesn't exist
            return {
                notFound: true
            }
        } else {
            // Throw the error again for us to see in logs
            throw e
        }
    }
}

/**
 * The specific parameters that we need to get from the URL
 */
interface IParams extends ParsedUrlQuery {
    id: string
}