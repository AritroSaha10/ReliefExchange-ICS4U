import Link from "next/link";
import { useEffect, useState } from "react";
import { SelectPicker } from "rsuite";
import { getCookie, hasCookie, setCookie } from 'cookies-next';
import Dropdown from "./Dropdown";

const languages = [
    { label: 'English', value: '/auto/en' },
    { label: 'French', value: '/auto/fr' },
    { label: 'Español', value: '/auto/es' },
    { label: 'हिन्दी', value: '/auto/hi' },
    { label: 'اردو', value: '/auto/ur' },
    { label: `Русский`, value: '/auto/ru' },
    { label: 'Polski', value: '/auto/pl' }
];

const languagesDropdown = languages.map(obj => ({
    ...obj,
    name: obj.label
}))

export default function Footer() {
    const [selectedLanguage, setSelectedLanguage] = useState(languages[0].value)

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
            };

            if (hasCookie('googtrans')) {
                setSelectedLanguage(getCookie('googtrans'))
            }
            else {
                setSelectedLanguage('/auto/en')
            }
        } catch (e) {
            console.error("Error while configuring auto-translate: ", e)
        }
    }, [])

    const langChange = (e) => {
        if (hasCookie('googtrans')) {
            setCookie('googtrans', decodeURI(e))
            setSelectedLanguage(e)
        }
        else {
            setCookie('googtrans', e)
            setSelectedLanguage(e)
        }
        window.location.reload()
    }

    return (
        <footer className="flex flex-col gap-3 p-4 bg-slate-800">
            <div className="flex justify-center items-center gap-3">
                <p className="text-gray-200 text-sm">Made by <a href="https://www.aritrosaha.ca/" className="text-blue-500 hover:text-blue-700 duration-200">Aritro Saha</a>, Joshua Chou, and Taha Khan</p>
            </div>

            <hr className="mx-16 md:mx-32 lg:mx-64 bg-slate-600 border-none h-px" />

            <div className="flex justify-center items-center gap-3">
                <Link href="/contacts" className="text-blue-500 hover:text-blue-700 duration-200">
                    Contact Us
                </Link>
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