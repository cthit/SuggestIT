import React, { Component, useState, useContext, useEffect } from 'react';
import './SuggestionBoard.css';
import { suggestions } from '../../../../redux/SuggestionStore';
import SuggestionItem from '../../../common/SuggestionItem';
import { updateSuggestions, deleteSuggestions } from '../../../../services/data.service';
import { DigitButton, DigitText, DigitDialogActions } from '@cthit/react-digit-components';
import { connect } from 'react-redux';
import { SuggestionsContext } from '../../../../common/suggestion-context';

const SuggestionBoardView = ({ dialogOpen }) => {
	const [ clearButton, setClearButton ] = useState(null);
	const [ suggestions, setSuggestions ] = useContext(SuggestionsContext);

	useEffect(() => {
		console.log('Updating suggestions');
		updateSuggestions(setSuggestions);
		return () => {};
	}, []);

	useEffect(
		() => {
			setClearButton(
				suggestions.length > 0 ? getClearButton() : <DigitText.Title text="The suggestion box is empty." />
			);
			return () => {};
		},
		[ suggestions ]
	);

	const getClearButton = () => (
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
	);

	const clearSuggestions = () => {
		deleteSuggestions(suggestions).then((res) => updateSuggestions(setSuggestions));
	};

	return (
		<div>
			<div className="grid">{clearButton}</div>
			<div className="grid">
				{suggestions.map((obj) => <SuggestionItem key={obj.id} suggestion={obj} ts={obj.timestamp} />)}
			</div>
		</div>
	);
};

const mapStateToProps = (state, ownProps) => ({});

const mapDispatchToProps = (dispatch) => ({
	dialogOpen: (dialogData) => dispatch(DigitDialogActions.digitDialogOpen(dialogData))
});

const SuggestionBoard = connect(mapStateToProps, mapDispatchToProps)(SuggestionBoardView);

export default SuggestionBoard;
