import React from 'react';
import { useNavigate } from 'react-router-dom';

function Home() {
  const navigate = useNavigate();

  const handleNavigateToCreateRoom = () => {
    navigate('/room/create')
  };

  return (
    <div>
      <h1>Welcome to the Home Page</h1>
      <p>This is the home page of your application.</p>
      <button onClick={handleNavigateToCreateRoom}>Create Room</button> 
    </div>
  );
}

export default Home;