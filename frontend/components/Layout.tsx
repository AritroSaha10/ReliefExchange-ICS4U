import Head from "next/head"

import Navbar from "./Navbar"
import Footer from "./Footer"

export default function Layout({ name, children }: { name: string, children: any }) {
    const title = `${name} | Relief Exchange`;
    const description = "Relief Exchange is a platform where generosity meets community. Learn more at our website.";

    return (
        <div className="flex flex-col min-h-screen bg-slate-900 overflow-hidden" key={name}>
            <Head>
                <title>{title}</title>
                <meta name="description" content={description} />

                <meta property="og:title" content={title} />
                <meta property="og:description" content={description} />
                <meta property="og:type" content="website" />

                <meta name="twitter:card" content="summary_large_image" />
                <meta property="twitter:title" content={title} />
                <meta property="twitter:description" content={description} />
            </Head>

            <Navbar />

            <div
                className="flex-grow flex flex-col"
            >
                {children}
            </div>

            <Footer />
        </div>
    )
}