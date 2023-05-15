import Link from "next/link";

export default function Footer() {
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
        </footer>
    );
}