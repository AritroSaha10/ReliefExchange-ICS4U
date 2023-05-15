import Head from "next/head"

import Navbar from "./Navbar"
import Footer from "./Footer"

export default function Layout({ name, children, noAnim }: { name: string, children: any, noAnim?: boolean }) {
    const title = `${name} | Relief Exchange`;
    const description = "Relief Exchange is a adsjlkkljasdkljadskljdas";
    const imageSrc = "CHANGE ME"

    return (
        <div className="flex flex-col min-h-screen bg-slate-900 overflow-hidden" key={name}>
            <Head>
                <title>{title}</title>
                <meta name="description" content={description} />

                <meta property="og:title" content={title} />
                <meta property="og:description" content={description} />
                <meta property="og:type" content="website" />
                <meta property="og:image" content={imageSrc} />
                <meta property="og:image:type" content="image/png" />
                <meta property="og:image:width" content="1111" />
                <meta property="og:image:height" content="1111" />

                <meta name="twitter:card" content="summary_large_image" />
                <meta name="twitter:creator" content="@YOUR_TWITTER" />
                <meta property="twitter:title" content={title} />
                <meta property="twitter:description" content={description} />
                <meta property="twitter:image:src" content={imageSrc} />
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