import Head from 'next/head';
import Link from 'next/link';
import Image from 'next/image';
import Layout from '../components/Layout';

import { useState, useEffect } from "react"
import { onAuthStateChanged } from 'firebase/auth';
import auth from '../lib/firebase/auth';

import HeroImage from "@media/hero.jpg";
import KeepDonateStock from "@media/keep-donate-stock.jpg"

export default function Home() {
  const [isSignedIn, setSignedIn] = useState(false);

  /**
   * Refresh auth data on change
   */
  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, user => {
      setSignedIn(user !== null);
    });

    return () => unsubscribe();
  }, []);

  return (
    <Layout name="Home">
      <header className="h-screen relative">
        <Image
          src={HeroImage}
          placeholder="blur"
          alt="Helping hand"
          fill
          quality={100}
          priority={true}
          className="object-cover object-center blur-sm"
        />

        <div className="relative z-1 h-screen bg-opacity-40 bg-black flex items-center justify-center">
          <div className="flex flex-col gap-4 mx-2 text-center">
            <h1 className="text-gray-500 font-bold text-4xl xs:text-5xl md:text-6xl">
              <span className="text-white">ReliefExchange</span>
            </h1>

            <h2 className="text-gray-500 font-semibold text-2xl xs:text-3xl md:text-4xl">
              <span className="text-gray-300">Where generosity meets community</span>
            </h2>

            <div className="flex flex-wrap gap-2 justify-center mt-4">
              {isSignedIn &&
                <Link className='px-4 py-2 bg-blue-500 hover:bg-blue-700 text-white text-xl active:bg-blue-800 duration-75 rounded-lg font-medium' href="/donations/create">
                  Donate
                </Link>
              }
              <Link className='px-4 py-2 bg-blue-500 hover:bg-blue-700 text-white text-xl active:bg-blue-800 duration-75 rounded-lg font-medium' href="/donations">
                View Recent Items
              </Link>
            </div>
          </div>
        </div>
      </header>

      <section className="flex p-10 flex-col items-center lg:flex-row lg:p-20 xl:px-40 items-left bg-transparent gap-6" id="about">
        <div
          className="flex flex-col items-center lg:items-start w-4/5 text-center lg:text-left mb-4 lg:mb-0"
        >
          <h1 className="text-white font-bold text-3xl md:text-4xl">
            Helping people <span className='text-blue-300'>help people</span>
          </h1>
          <p className="mt-4 w-full md:w-3/4 text-lg text-gray-200">
            We&#39;re a student-led organization that&apos;s focused on
            getting help to people. Here&apos;s the gist: you have some spare
            stuff. Whether its a fish tank or some exercise bands, you just
            don&apos;t need it now. Instead of throwing it away, why not give it to
            someone else who might not be able to regularly buy it, without you
            having to going through all the hassle of selling it?
          </p>

          <Link href="/donations" className="group mt-6 bg-blue-300 text-black font-semibold py-2 px-4 rounded-lg text-lg hover:bg-blue-400 duration-75">
            View Recent Items <span className="group-hover:ml-1 duration-150 transition-all">→</span>
          </Link>
        </div>

        <div
          className="flex p-0 m-0 w-1/2"
        >
          <Image
            src={KeepDonateStock}
            alt="Donation stock photo"
            className="object-cover object-center rounded-xl"
            width={900}
            height={540}
            placeholder="blur"
            loading="lazy"
          />
        </div>
      </section>
    </Layout>
  )
}
