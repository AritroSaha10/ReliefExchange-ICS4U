import Link from "next/link";
import Image from "next/image";
import styles from "../..styles/donation.css";
export default function donation(donation)
{
return(
  <>
    <div className={styles.container}>
        <Image src={donation.src}/> 
    </div>
              <Link href="/donations">Back to Donations</Link>
              </>
)
}