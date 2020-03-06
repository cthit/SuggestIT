import React, { useState, useContext, useEffect, useCallback } from 'react';
import './SuggestionBoard.css';
import SuggestionItem from '../../../common/SuggestionItem';
import { updateSuggestions, deleteSuggestions } from '../../../../services/data.service';
import { DigitButton, DigitText, useDigitDialog } from '@cthit/react-digit-components';
import { SuggestionsContext } from '../../../../common/suggestion-context';

const SuggestionBoard = () => {
	const [ dialogOpen ] = useDigitDialog();
	const [ clearButton, setClearButton ] = useState(null);
	const [ suggestions, setSuggestions ] = useContext(SuggestionsContext);

	useEffect(
		() => {
			console.log('Updating suggestions');
			updateSuggestions(setSuggestions);
			return () => {};
		},
		[ setSuggestions ]
	);

	const clearSuggestions = useCallback(
		() => deleteSuggestions(suggestions).then((res) => updateSuggestions(setSuggestions)),
		[ suggestions, setSuggestions ]
	);

	const getClearButton = useCallback(
		() => (
			<DigitButton
				className="clear-button"
				text="Clear suggestions"
				primary
				raised
				onClick={() => {
					dialogOpen({
						title: 'Are you sure',
						description: `Are you sure you want to delete ${suggestions.length} suggestions?`,
						cancelButtonText: 'No',
						confirmButtonText: 'Yes',
						onConfirm: () => clearSuggestions()
					});
				}}
			/>
		),
		[ suggestions, clearSuggestions, dialogOpen ]
	);

	useEffect(
		() => {
			setClearButton(
				suggestions.length > 0 ? getClearButton() : <DigitText.Title text="The suggestion box is empty." />
			);
			return () => {};
		},
		[ getClearButton, setClearButton, suggestions ]
	);

	return (
		<div>
			<div className="grid">{clearButton}</div>
			<div className="grid">
				{suggestions.map((obj) => <SuggestionItem key={obj.id} suggestion={obj} ts={obj.timestamp} />)}
			</div>
		</div>
	);
};

export default SuggestionBoard;
