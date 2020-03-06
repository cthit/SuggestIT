import React, { Component, useState, useEffect } from 'react';

import './home.css';
import Prompt from './views/Prompt/prompt.view';
import SuggestionBoard from './views/SuggestionBoard/SuggestionBoard';
import { checkLogin } from '../../services/data.service';

const Home = () => {
	const [ suggestionboard, setSuggestionBoard ] = useState(<div />);

	useEffect(() => {
		checkLogin().then(() => setSuggestionBoard(<SuggestionBoard />));
		return () => {};
	}, []);

	return (
		<div className="main">
			<Prompt />
			{suggestionboard}
		</div>
	);
};

export default Home;
