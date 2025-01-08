import React from "react";
import { useNavigate } from "react-router-dom";

const CreateRoom = () => {
    const navigate = useNavigate();

    const create = async (e) => {
        e.preventDefault();

        try {
            const resp = await fetch('http://localhost:8080/room/create', {
                method: 'GET', // Ensure the correct method is used
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            if (!resp.ok) {
                throw new Error('Network response was not ok');
            }
           
            const { roomID } = await resp.json();
            console.log('Room ID:', roomID);
            // once we get roomid, we need to redirect to the room using the roomid so we use below
            navigate(`/room/${roomID}`);
        } catch (error) {
            console.error('There was a problem with the fetch operation:', error);
        }
    };

    return (
        <div>
            <h1>Create a Room</h1>
            <button onClick={create}>Create Room</button>
        </div>
    );
};

export default CreateRoom;