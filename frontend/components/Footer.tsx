/**
 * @file Footer component definition.
 * 
 * This file contains the implementation of the Footer component, which represents the footer of
 * the website. It includes functionality for translating the page using the Google Translate API.
 * The component imports functions from 'cookies-next' for managing cookies and the 'Dropdown'
 * component for selecting a language. The supported languages for auto-translation are defined in
 * the 'languages' array. The component uses the 'useState' and 'useEffect' hooks from React.
 * 
 * @author Aritro Saha
 * @cite Installation - Tailwind CSS, https://tailwindcss.com/docs/installation. 
 * @cite “Headless ui,” Headless UI, https://headlessui.com/react/listbox. 
 * @cite “React Icons,” React Icons, https://react-icons.github.io/react-icons/.
 */


import { useEffect, useState } from "react";
import { deleteCookie, getCookie, hasCookie, setCookie } from 'cookies-next';
import Dropdown from "./Dropdown";

// All the languages supported for auto-translation
const languages = [
    { label: 'English', value: '/auto/en' },
    { label: 'French', value: '/auto/fr' },
    { label: 'Español', value: '/auto/es' },
    { label: 'हिन्दी', value: '/auto/hi' },
    { label: 'اردو', value: '/auto/ur' },
    { label: `Русский`, value: '/auto/ru' },
    { label: 'Polski', value: '/auto/pl' }
];

// Language list, converted into a format suitable for the Dropdown class
const languagesDropdown = languages.map(obj => ({
    ...obj,
    name: obj.label
}))

/**
 * The footer of the website, with the ability to translate the page.
 */
export default function Footer() {
    // State to remember the selected language
    const [selectedLanguage, setSelectedLanguage] = useState(languages[0].value)

    /**
     * Set up auto-translate functionality
     */
    useEffect(() => {
        try {
            // @ts-ignore All of this is imported by the extra script
            window.googleTranslateElementInit = () => {
                // @ts-ignore
                new window.google.translate.TranslateElement({
                    pageLanguage: 'auto',
                    autoDisplay: false,
                    includedLanguages: "ru,en,pl,fr,es,hi,ur", // If you remove it, by default all google supported language will be included
                    // @ts-ignore
                    layout: google.translate.TranslateElement.InlineLayout.SIMPLE
                },
                    'google_translate_element');

                // Use the cookie value if the user wanted to translate it to another language before
                if (hasCookie('googtrans')) {
                    setSelectedLanguage(getCookie('googtrans').toString())
                }
                else {
                    // Auto-translate to english on default
                    setSelectedLanguage('/auto/en')
                }
            };
        } catch (e) {
            // Silently log error
            console.error("Error while configuring auto-translate: ", e)
        }
    }, [])

    /**
     * Change the language to translate to
     * @param lang The language code to translate to
     */
    const langChange = (lang) => {
        // Delete any cookies made by our applications or Google Translate
        deleteCookie("googtrans", { path: "", domain: `.${location.hostname}`})
        deleteCookie("googtrans", { path: "", domain: `${location.hostname}`})

        // Set new cookie and update state
        setCookie('googtrans', lang)
        setSelectedLanguage(lang)

        // Reload the page for the new language to show
        window.location.reload()
    }

    return (
        <footer className="flex flex-col gap-3 p-4 bg-slate-800">
            <div className="flex justify-center items-center gap-3">
                <p className="text-gray-200 text-sm">Made by <a href="https://www.aritrosaha.ca/" className="text-blue-500 hover:text-blue-700 duration-200">Aritro Saha</a>, Joshua Chou, and Taha Khan</p>
            </div>

            <hr className="mx-16 md:mx-32 lg:mx-64 bg-slate-600 border-none h-px" />

            <div className="flex justify-center items-center gap-3">
                <a href="https://docs.google.com/document/d/1SwvbGomqzTCoZS3yOiGqkT-3oLW5rO9EociF8mV30MM/edit" target="_blank" rel="noreferrer" className="text-blue-500 hover:text-blue-700 duration-200">
                    Quickstart Guide
                </a>
            </div>

            <div id="google_translate_element" style={{ width: '0px', height: '0px', position: 'absolute', left: '50%', zIndex: -99999 }}></div>

            <div className="flex flex-col items-center justify-center w-full">

                <Dropdown
                    title="Select a language"
                    selectedItem={languagesDropdown.filter(obj => obj.value === selectedLanguage)[0]}
                    setSelectedItem={(langObj) => langChange(langObj.value)}
                    options={languagesDropdown}
                    openOverlap={false}
                />
            </div>
        </footer>
    );
}