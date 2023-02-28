import React, { useEffect, useState } from 'react';

function Details(props) {
  const [donation, setDonation] = useState(null);

  useEffect(() => {
    // Send a request to the server to retrieve the details of the selected donation.
    fetch(`/donations/${props.id}`)
      .then(response => response.json())
      .then(data => setDonation(data));
  }, [props.id]); //will rerender if props.id (second argument) changes 

  if (!donation) {
    return <div>Loading...</div>;
  }

  return ( 
    <div>
      <h2>{donation.description}</h2>
      <p>{donation.location}</p>
    </div>
  );
}

export default Details;