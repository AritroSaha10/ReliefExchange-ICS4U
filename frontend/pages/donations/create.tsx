import Layout from "@components/Layout";
import { useState, useEffect, useRef } from "react";
import { getIdToken, onAuthStateChanged, User } from "firebase/auth";
import auth from "lib/firebase/auth";
import { useRouter } from "next/router";
import Multiselect from 'multiselect-react-dropdown';
import ReCAPTCHA from "react-google-recaptcha"
import axios from "axios";
import React from "react";
import { AiOutlineCloudUpload } from "react-icons/ai"
import { BsImage } from "react-icons/bs"
import dynamic from "next/dynamic";
import * as commands from "@uiw/react-md-editor/lib/commands";
import { getDownloadURL, ref, uploadBytes } from "firebase/storage";
import storage from "lib/firebase/storage";
import allTags from "lib/tag-types";

import "@uiw/react-md-editor/markdown-editor.css";
import "@uiw/react-markdown-preview/markdown.css";
import convertBackendRouteToURL from "lib/convertBackendRouteToURL";

const MDEditor = dynamic(
    () => import("@uiw/react-md-editor").then((mod) => mod.default),
    { ssr: false }
);

export default function CreateDonation() {
    const router = useRouter();
    const captchaRef = useRef(null);

    const [loadingAuth, setLoadingAuth] = useState(true);
    const [user, setUser] = useState<User>(null);
    const [signedIn, setSignedIn] = useState<boolean>(false);
    const [tagsSelected, setTagsSelected] = useState([]);
    const [descriptionMD, setDescriptionMD] = useState("**Hello world!!!**");
    const [submitting, setSubmitting] = useState(false);
    const [featuredImage, setFeaturedImage] = useState<FileList | []>([]);

    const RECAPTCHA_SITE_KEY = process.env.NEXT_PUBLIC_RECAPTCHA_SITE_KEY;

    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, newUser => {
            if (newUser && Object.keys(newUser).length !== 0) {
                // Set user data
                setUser(newUser);
                setSignedIn(true);
            } else {
                setUser(null);
                setSignedIn(false);
                alert("You need to be signed in to access this page. Redirecting...");
                router.push("/");
            }

            setLoadingAuth(false);
        });

        return () => unsubscribe();
    }, [router]);

    const onSubmit: React.FormEventHandler<HTMLFormElement> = async (e) => {
        const formData = Object.fromEntries((new FormData(e.currentTarget)).entries());
        e.preventDefault();

        setSubmitting(true);

        // Confirming the CAPTCHA
        const token = captchaRef.current.getValue();
        captchaRef.current.reset();
        if (!token) {
            alert("Please complete the CAPTCHA and try again.");
            return;
        }
        try {
            const res = await axios.post(convertBackendRouteToURL(`/confirmCAPTCHA?token=${token}`))
            if (!res.data.human) throw "User was detected to be a bot by ReCAPTCHA."
        } catch (e) {
            alert("Something went wrong. Please try again.");
            console.error(e);
            setSubmitting(false);
            return
        }

        // CAPTCHA confirmed, now upload the image to Firebase Storage
        let imgLink = ""
        if (featuredImage.length !== 0) {
            const imgRef = ref(storage, `donations/${crypto.randomUUID()}.jpg`);

            try {
                const imgSnapshot = await uploadBytes(imgRef, featuredImage[0]);
                imgLink = await getDownloadURL(imgSnapshot.ref);
            } catch (e) {
                alert("Something went wrong while uploading your image. Please try again, and make sure that your image is <=10MB.");
                console.error(e);
                setSubmitting(false);
                return
            }
        }

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
            "img": imgLink,
            "tags": tagsSelected.map(obj => obj.name),
            "creation_timestamp": nowUTC.toISOString(),
            "ownerID": user.uid
        };

        // Send the prep'd data to our endpoint
        try {
            const apiRes = await axios.post(convertBackendRouteToURL("/donations/new"), {
                data: donationData,
                token: idToken
            });
            alert("Your donation was successfully submitted! Redirecting you to its page...");
            router.push(`/donations/${apiRes.data}`);
        } catch (e) {
            alert("Something went wrong while submitted your donation. Please try again.");
            console.error(e);
            setSubmitting(false);
            return;
        }

        setSubmitting(false);
    }

    return (
        <Layout name="Make a Donation">
            <div className="flex flex-col gap-4 p-10 flex-grow min-w-screen">
                <div className="flex flex-col gap-2 mb-4">
                    <h1 className="text-4xl text-white font-bold text-center">Make Donation Offer</h1>
                    <p className="text-md text-gray-200 text-center">Fill out this form to offer a donation to others!</p>
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
                                        maxLength={100}
                                    />
                                </div>
                            </div>

                            <div className="flex flex-col items-center">
                                <h3 className="text-white text-2xl font-medium mb-2 text-center lg:text-left">Product Description: (Max. 1000 characters) <span className="text-red-500"> *</span></h3>
                                <div className="flex flex-col gap-4 w-full items-center" data-color-mode="light">
                                    <MDEditor
                                        value={descriptionMD}
                                        onChange={setDescriptionMD}
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
                                        maxLength={100}
                                    />
                                </div>
                            </div>

                            <div className="flex flex-col items-center">
                                <h3 className="text-white text-2xl font-medium mb-2 text-center">Product Photo:</h3>
                                <div className="flex flex-col gap-4 items-center w-full lg:w-1/2">
                                    <label className={`${featuredImage.length ? "bg-slate-500" : "bg-gray-200"} text-white font-semibold py-2 px-4 rounded-xl border-2 border-slate-200 hover:border-slate-400 duration-300 lg:w-3/4 shadow cursor-pointer`}>
                                        <div className="flex flex-col items-center justify-center">
                                            {featuredImage.length == 0 ? (
                                                <>
                                                    <AiOutlineCloudUpload className="text-6xl text-slate-500" />
                                                    <h2 className="text-xl text-slate-500 font-semibold text-center">Upload An Image</h2>
                                                </>
                                            ) : (
                                                <>
                                                    <BsImage className="text-6xl text-white" />
                                                    <h2 className="text-xl text-white font-semibold text-center">Attached {featuredImage[0].name}</h2>
                                                </>
                                            )
                                            }
                                        </div>
                                        <input
                                            type="file"
                                            className="absolute w-px h-px p-0 -m-px overflow-hidden border-0"
                                            style={{ clip: "rect(0, 0, 0, 0)" }}
                                            accept="image/png, image/gif, image/jpeg"
                                            name="featuredImage"
                                            onChange={
                                                (e) => setFeaturedImage(e.target.files)
                                            }
                                        />
                                    </label>
                                </div>
                                {featuredImage.length ?
                                    <button className="text-md text-red-500" type="button" onClick={() => setFeaturedImage([])}>Remove</button>
                                    : ""
                                }
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