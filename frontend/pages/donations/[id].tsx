import axios from "axios";
import Image from "next/image";

import { GetStaticPaths, GetStaticProps } from 'next'
import { ParsedUrlQuery } from 'querystring'
import Donation from "lib/types/donation";
import Layout from "@components/Layout";
import { ReactMarkdown } from "react-markdown/lib/react-markdown";
import UserData from "lib/types/userData";

import { BiLeftArrowAlt } from "react-icons/bi"
import Link from "next/link";

interface IParams extends ParsedUrlQuery {
    id: string
}

interface RawDonation {
    id: string,
    title: string,
    description: string, // markdown
    location: string,
    img: string, // direct src to firebase image
    creation_timestamp: string,
    tags: string[],
    owner_id: string
}

interface DonationWithUserData extends Donation {
    owner: UserData
}

export const getStaticPaths: GetStaticPaths = async () => {
    const rawDonations: RawDonation[] = (await axios.get("http://localhost:8080/donations/list")).data
    const donations: Donation[] = rawDonations.map(rawDonation => ({
        ...rawDonation,
        creation_timestamp: new Date(rawDonation.creation_timestamp)
    }))

    const arr: string[] = donations.map(donation => donation.id)

    return {
        paths: arr.map((id) => {
            return {
                params: { id },
            }
        }),
        fallback: "blocking"
    }
}

export const getStaticProps: GetStaticProps = async (context) => {
    const { id } = context.params as IParams

    try {
        const rawDonation: Donation = (await axios.get(`http://localhost:8080/donations/${id}`)).data
        const donation: DonationWithUserData = {
            ...rawDonation,
            owner: (await axios.get(`http://localhost:8080/users/${rawDonation.owner_id}`)).data
        }

        const props = { donation }
        return { props }
    } catch (e) {
        console.error(e);
        return {
            notFound: true
        }
    }
}

export default function DonationSpecificPage({ donation }: { donation: DonationWithUserData }) {
    return (
        <Layout name={donation.title}>
            <div className="flex flex-col lg:items-center lg:justify-center flex-grow">
                <div className="flex flex-col lg:flex-row px-10 py-4 lg:px-20 lg:py-4 xl:px-60 xl:py-10 items-center">
                    <div className="flex flex-col gap-2 items-center lg:w-1/2">
                        <Link href="/donations" className="flex items-center text-blue-500 hover:text-blue-600 active:text-blue-700 duration-150 lg:self-start mb-2">
                            <BiLeftArrowAlt />
                            Back to Donations
                        </Link>

                        {donation.img ? <Image src={donation.img} alt="Featured image" height={500} width={500} className="rounded-md object-cover object-center" /> : <></>}
                    </div>

                    <div className="lg:ml-5 flex flex-col items-center lg:items-start">
                        <h1 className="text-white text-4xl font-semibold text-center break-all">{donation.title}</h1>
                        <h3 className="text-gray-300 text-md text-center break-all">Available in "{donation.location}"</h3>

                        <br />

                        {/* This makes sure that all external links in the description open another tab */}
                        <base target="_blank" />

                        <ReactMarkdown skipHtml={true} className="text-white break-all">{donation.description}</ReactMarkdown>

                        <br />

                        <a className="py-2 px-4 bg-blue-500 font-semibold text-center text-white rounded-lg hover:bg-blue-600 duration-75" href={`mailto:${donation.owner.email}`}>Contact {donation.owner.display_name} for More Info</a>
                    </div>
                </div>
            </div>
        </Layout>
    )
}
