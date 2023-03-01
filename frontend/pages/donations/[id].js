import Link from "next/link";
import Image from "next/image";
 import Layout from "../../components/Layout" 
export default function donation(donation)
{
return(
  <Layout>

  <Link href="/donations">Back to Donations</Link>
        <Image src={donation.src}/> 
              </Layout>
)
}
export async function getProps(context)
{
  const {id}=context.query; //make id avalible 
  const res =await fetch(`http://localhost:4000/donations/${id}`)
  const donation=await res.json()
  return{
    props:donation,
  }
}