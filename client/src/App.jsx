import React, { useState } from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import CreateRoom from './components/CreateRoom';
import Home from './components/Home';
import Room from './components/Room';
import GetRoom from './components/GetRoom';

function App() {
  const [count, setCount] = useState(0);

  return (
    <div className="App">
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/room/create" element={<CreateRoom />} />
          <Route path="/room/:roomID" element={<GetRoom />} />
          <Route path="/room/:roomID/join" element={<Room />} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;