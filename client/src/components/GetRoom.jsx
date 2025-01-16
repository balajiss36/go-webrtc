import React, { useEffect } from 'react';
import { useParams,useNavigate } from 'react-router-dom';


function GetRoom() {
    const navigate = useNavigate();
    const { roomID } = useParams();
  
    const get = async (e) => {
        e.preventDefault();

        try {
            const resp = await fetch(`http://localhost:8080/room/${roomID}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            if (!resp.ok) {
                throw new Error('Network response was not ok');
            }

            // once we get roomid, we need to redirect to the room using the roomid so we use below
            navigate(`/room/${roomID}/join`);
        } catch (error) {
            console.error('There was a problem with the fetch operation:', error);
        }
    };
  
    return (
        <div>
        <h1>Join a Room</h1>
        <button onClick={get}>Join Room</button>
    </div>
    );
  }

export default GetRoom;