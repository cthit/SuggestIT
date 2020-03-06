import React, { useState, useContext } from 'react';
import { DigitNavLink, DigitText, DigitIconButton, DigitLayout, useDigitToast } from '@cthit/react-digit-components';
import './suggestionitem.style.css';
import { translateTimestamp } from '../methods';
import { Clear, ExpandMore, ExpandLess } from '@material-ui/icons';
import { updateSuggestions, deleteSuggestion, addSuggestion } from '../../../services/data.service';
import { SuggestionsContext } from '../../../common/suggestion-context';

const SuggestionItem = ({ suggestion, ts, ...props }) => {
	const [ text, dispatchToggle ] = useState(null);
	const [ toastOpen ] = useDigitToast({ duration: 5000 });
	const [ , setSuggestions ] = useContext(SuggestionsContext);

	const toggleText = () => {
		dispatchToggle(text ? null : <DigitText.Subtitle className="suggestionText" text={suggestion.text} />);
	};

	const deleteWithToast = () => {
		toastOpen({
			text: 'The suggestion has been deleted',
			actionText: 'Undo',
			actionHandler: () => {
				addSuggestion(suggestion).then((res) => updateSuggestions(setSuggestions));
			}
		});
		deleteSuggestion(suggestion.id).then((res) => updateSuggestions(setSuggestions));
	};

	return (
		<div className="card">
			<DigitLayout.Grid>
				<DigitLayout.Row>
					<DigitNavLink text={suggestion.title} link={'/suggestion/' + suggestion.id} />
					<div className="float-right">
						<DigitIconButton icon={Clear} onClick={() => deleteWithToast()} />
					</div>
				</DigitLayout.Row>
				<DigitLayout.Row>
					<DigitText.Subtitle2 text={suggestion.author} className="authorLabel" />
					<div className="float-right">
						<DigitIconButton icon={text ? ExpandLess : ExpandMore} onClick={() => toggleText()} />
					</div>
				</DigitLayout.Row>
				<DigitLayout.Row>
					<DigitText.Subtitle2 text={translateTimestamp(ts)} className="timeStampLabel" />
				</DigitLayout.Row>
				{text}
			</DigitLayout.Grid>
		</div>
	);
};

export default SuggestionItem;
