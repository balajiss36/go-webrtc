import React, { useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';

const Room = () => {
  const { roomID } = useParams();
  const ws = useRef(null);

  const userVideo = useRef();
  const partnerVideo = useRef();
  const userStream = useRef();
  const peerRef = useRef();
  const websocketRef = useRef();
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
        ws.current = new WebSocket(`ws://localhost:8080/room/${roomID}`);
        ws.current.addEventListener('open', () => {
          console.log('Connected to the server');
          ws.current.send(JSON.stringify({ join: "true" }));
        });
    
        ws.current.addEventListener('message', async (event) => {
            const data = JSON.parse(event.data);
            console.log('Received data:', data);

            if (message.join){
                    callUser();
            }
            // when we find a ICE candidate
            if (message.iceCandidate){
                console.log('Received ICE candidate:', message.iceCandidate);
                try {
                    peerRef.current.addIceCandidate(message.iceCandidate);
                } catch (error) {
                    console.log('Error receiving ICE candidate:', error);
                }
            }

            // to accept the offer
            if (message.offer){
                handleOffer(message.offer);
            }

            if (message.answer){
                console.log('Received answer:', message.answer);
                peerRef.current.setRemoteDescription(new RTCSessionDescription(message.answer));
            }
        });
    
        return () => {
          if (ws.current) {
            ws.current.close();
          }
        };
      }, [roomID]);
    })

    const handleOffer = async (offer) => {
        console.log('Received offer, Creating Answer', offer);
        peerRef.current = createPeer();

        // gets the offer and sets it as the remote description
        const desc = new RTCSessionDescription(offer);
        await peerRef.current.setRemoteDescription(desc);

        userStream.current.getTracks().forEach((track) => {
            peerRef.current.addTrack(track, userStream.current)
            });

            // peer will set answer as local description and send it to the remote peer
        const answer = await peerRef.current.createAnswer();
        await peerRef.current.setLocalDescription(answer);
    };

    const callUser = () => {
        console.log('Calling user');
        peerRef.current = createPeer();
        // add all the media tracks (audio and video) from the user's media stream to the peer connection. This allows the media to be sent to the remote peer.
        userStream.current.getTracks().forEach((track) => {
            peerRef.current.addTrack(track, userStream.current)
            });
    };

    const createPeer = async ()=> {
        const peer = new RTCPeerConnection({
            iceServers: [
                {
                    urls: 'stun:stun.l.goolge.com:19302',
                },
            ],
        });

        peer.onnegotiationneeded = () => handleNegotiationNeededEvent(peer);
        peer.onicecandidate = handleICECandidateEvent;
        peer.ontrack = handleTrackEvent;
        return peer
    };

    const handleNegotiationNeededEvent = async (peer) => {
        console.log('Creating offer');
        try {
            const offer = await peer.createOffer();
            await peer.setLocalDescription(offer);
            console.log('Sending offer:', offer);
            ws.current.send(
                JSON.stringify({
                    offer: offer,
                })
            );
        }
        catch (error) {
            console.error('Error creating offer:', error);
        }
    };
    // whenever the connection finds a ICE candidate, it sends it to the remote peer using the WebSocket connection via signallint server.
    const handleICECandidateEvent = (event) => {
        if (event.candidate) {
            console.log('Sending ICE candidate:', event.candidate);
            ws.current.send(
                JSON.stringify({
                    ice: event.candidate,
                })
            );
        }
    };

    const handleTrackEvent = (event) => {
        console.log('Received a track:', event);
        partnerVideo.current.srcObject = event.streams[0];
    };

  return (
    <div>
      <video autoPlay controls={true} ref={userVideo}></video>
      <video autoPlay controls={true} ref={partnerVideo}></video>
    </div>
  );
};

export default Room;