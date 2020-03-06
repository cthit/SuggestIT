import React, { useEffect } from 'react';
import Cookies from 'universal-cookie';
const cookies = new Cookies();

const Callback = ({ location }) => {
	useEffect(
		() => {
			let params = new URLSearchParams(location.search);
			cookies.set('AUTH_TOKEN', params.get('token'));
			window.location.replace('/');
			return () => {};
		},
		[ location.search ]
	);

	return <div />;
};

export default Callback;
