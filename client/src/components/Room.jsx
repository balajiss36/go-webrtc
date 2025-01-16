import React, { useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';

const Room = () => {
  const { roomID } = useParams();
  const ws = useRef(null);
  const userVideo = useRef();
  const partnerVideo = useRef();
  const userStream = useRef();
  const openCamera = async () => {
    try {
      const devices = await navigator.mediaDevices.enumerateDevices();
      const allDevices = Array.isArray(devices) ? devices : Array.from(devices);
      const videoDevices = allDevices.filter(device => device.kind === 'videoinput');
      console.log('Video devices:', videoDevices);
  
      if (videoDevices.length === 0) {
        throw new Error('No video input devices found.');
      }
  
      const camera = videoDevices[0];
  
      const constraints = {
        audio: true,
        video: {
          deviceId: camera.deviceId,
        },
      };
  
      // gets the user media devices and returns the stream below
      return await navigator.mediaDevices.getUserMedia(constraints);
    } catch (error) {
      console.error('Error accessing media devices.', error);
    }
  };

  useEffect(() => {
    openCamera().then((stream) => {
      userVideo.current.srcObject = stream;
      userStream.current = stream;
    // Establish WebSocket connection
    ws.current = new WebSocket(`ws://localhost:8080/room/${roomID}/join`);

    ws.current.addEventListener('open', () => {
      console.log('Connected to the server');
      ws.current.send(JSON.stringify({ join: "true" }));
    });

    ws.current.addEventListener('message', (event) => {
      const data = JSON.parse(event.data);
      console.log('Received data:', data);

      // Handle different types of messages from the server
      if (data.join) {
        console.log('User joined the room');
      }

      if (data.iceCandidate) {
        console.log('Received ICE candidate');
      }

      // Add more message handling as needed
    });

    ws.current.addEventListener('close', () => {
      console.log('Disconnected from the server');
    });

    ws.current.addEventListener('error', (error) => {
      console.error('WebSocket error:', error);
    });

    // Cleanup on component unmount
    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [roomID]);
});

  return (
    <div>
      <video autoPlay controls={true} ref={userVideo}></video>
      <video autoPlay controls={true} ref={partnerVideo}></video>
    </div>
  );
};

export default Room;