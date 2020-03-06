import React, { useState, useEffect } from 'react';
import { DigitText } from '@cthit/react-digit-components';
import './suggestion.style.css';
import { translateTimestamp } from '../common/methods';
import { getSuggestion } from '../../services/data.service';

const Suggestion = ({ match }) => {
	const [ suggestion, setSuggestion ] = useState({ title: '', author: '', timestamp: '', text: '' });

	useEffect(
		() => {
			getSuggestion(match.params.id)
				.then((res) => {
					if (!res) {
						//Redirects to / if no suggestion could be found
						window.location.href = '/';
						return;
					}
					setSuggestion(res.data);
				})
				.catch((err) => console.log('Failed to get suggestion'));
			return () => {};
		},
		[ match.params.id ]
	);

	return (
		<div className="suggestionCard main">
			<DigitText.Heading6 text={suggestion.title} />
			<DigitText.Subtitle2 className="grayText" text={'Posted by: ' + suggestion.author} />
			<DigitText.Subtitle2 className="grayText" text={translateTimestamp(suggestion.timestamp)} />
			<DigitText.Text text={suggestion.text} />
		</div>
	);
};

export default Suggestion;
