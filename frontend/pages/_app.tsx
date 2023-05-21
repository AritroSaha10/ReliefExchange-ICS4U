import Script from "next/script";
import "../styles/globals.css";

export default function App({ Component, pageProps }) {
    return (
        <>
            <Component {...pageProps} />
            <Script src="//translate.google.com/translate_a/element.js?cb=googleTranslateElementInit" defer />
        </>
    );
}