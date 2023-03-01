import Link from "next/link";
import Image from "next/image";
import Head from "next/head"
//.. is move out of folder . is move out of file
import Layout from "../../components/Layout" 

export default function donations({donations}) {

    return(
        <Layout>
        <Head> <title>Donations</title></Head>
           <h1>Donations</h1>
        <h2>  <Link href="/">Back to home page</Link></h2>
<ul>
  {donations.map((donation)=>(
  <li key={donation.id}>
    <Link href={`/donations/${donation.id}`}><a>{donation.name}</a></Link>
    </li>
  ))}
</ul>
  {/* <Image src="/images/profile.jpg" height="144" width="120" alt="profile"/> */}

    </Layout>
    ) 
  }
  export async function getProps()//server side 
  {
const res=await fetch("https://localhost:4000/donations") //get donations from server side 
const donations=await res.json();
return {
  props: {donations},
};
  }