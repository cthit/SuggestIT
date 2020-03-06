import React, { useState, useContext } from 'react';
import {
	DigitText,
	DigitTextField,
	DigitTextArea,
	DigitButton,
	DigitSwitch,
	useDigitCustomDialog,
	useDigitToast
} from '@cthit/react-digit-components';
import { addSuggestion, updateSuggestions } from '../../../../services/data.service';
import './prompt.style.css';
import { LinkText } from './promtp.style';
import { SuggestionsContext } from '../../../../common/suggestion-context';

const title_error_message = 'The title is not filled in';
const description_error_message = 'The description is not filled in';

const Prompt = () => {
	const [ openDialog, , ] = useDigitCustomDialog();
	const [ toastOpen ] = useDigitToast({ duration: 5000 });

	const [ text, setText ] = useState('');
	const [ title, setTitle ] = useState('');
	const [ author, setAuthor ] = useState('');

	const [ errors, setErrors ] = useState({
		title_error: false,
		description_error: false
	});
	const [ anonymous_author, setAnonymousAuthor ] = useState(false);
	const [ , setSuggestions ] = useContext(SuggestionsContext);

	const sendNewSuggestion = (suggestion) => {
		addSuggestion(suggestion).then((res) => {
			updateSuggestions(setSuggestions);
			setTitle('');
			setAuthor('');
			setText('');
			toastOpen({
				text: 'Thank you! The suggestion has been sent to P.R.I.T.'
			});
		});
	};

	return (
		<div className="prompt">
			<div className="innerPrompt">
				<DigitText.Heading6 text="New Suggestion" />
				<DigitTextField
					error={errors.title_error}
					errorMessage={title_error_message}
					onChange={(e) => setTitle(e.target.value)}
					value={title}
					upperLabel="Title"
				/>
				<DigitTextArea
					error={errors.description_error}
					errorMessage={description_error_message}
					onChange={(e) => setText(e.target.value)}
					value={text}
					upperLabel="Suggestion"
					rows={5}
					rowsMax={10}
				/>
				<DigitTextField
					onChange={(e) => setAuthor(e.target.value)}
					value={author}
					disabled={anonymous_author}
					upperLabel="Name/Nick"
				/>
				<div>
					<DigitText.Text text={'By clicking Send I agree to'} />
					<LinkText
						onClick={() => {
							openDialog({
								renderButtons: () => null,
								renderMain: () => (
									<div>
										<DigitText.Title text={'GDPR agreement'} />
										<DigitText.Text
											text={
												"The data is collected by the Student Division Information Technology with the aim of knowing how to improve Hubben 2.1. The data that will be saved is title, the suggestion, the time it was sent and name if not anonymous. You have the right to withdraw your consent to the management, request to have all data that the organization has about you and to complain to the Data Inspectorate in case of dissatisfaction. The data will be kept until P.R.I.T. has solved the issue or deemed it invalid and will only be handled by the Student Division Information Technology, ultimately responsible for the data management of the section is William Levén, who can be reached at ordforande@chalmers.it. If you have any questions, please contact P.R.I.T. prit@chalmers.it or the section's data protection representative dpo@chalmers.it."
											}
										/>
									</div>
								)
							});
						}}
						text={'this GDPR agreement'}
					/>
				</div>
				<DigitSwitch
					value={anonymous_author}
					label="Anonymous"
					primary
					onChange={(e) => {
						setAnonymousAuthor(e.target.checked);
					}}
				/>
				<DigitButton
					text="Send"
					primary
					raised
					onClick={() => {
						setErrors({
							title_error: title === '',
							description_error: text === ''
						});

						if (title === '' || text === '') return;

						sendNewSuggestion({
							title: title,
							text: text,
							author: author === '' || anonymous_author ? 'Anonymous' : author
						});
					}}
				/>
			</div>
		</div>
	);
};

export default Prompt;
