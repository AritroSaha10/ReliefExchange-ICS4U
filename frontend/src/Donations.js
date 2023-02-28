import React, { useState, useEffect } from 'react';

function Donations() {
  const [donations, setDonations] = useState([]);
  useEffect(() => {
    async function fetchDonations() {
      const response = await fetch('http://localhost:4000/donations'); //using fetch instead of axios 
      const data = await response.json();
      setDonations(data);
    }
    fetchDonations();
  }, []);
  
  return (
    <div>
      <h1>Donations</h1>
      <ul>
        <img src={donation.src}/>
        {donations.map((donation) => (
          <li key={donation.id}>
            <p>{donation.description}</p>
            <p>{donation.location}</p>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default Donations;