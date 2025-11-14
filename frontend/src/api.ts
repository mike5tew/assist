// api.ts

import axios from 'axios';

export async function post(url: string, data: any) {
	try {
		const response = await axios.post(url, data);
		return response.data;
	} catch (error) {
		if (axios.isAxiosError(error)) {
			console.error(`API Response Error: ${error.message}`);
			console.error(`POST request failed for: ${url}`);
			console.error(error.response?.data || "No response data");
		} else {
			console.error(`Unexpected Error: ${error}`);
		}
		throw error; // Re-throw the error for further handling
	}
}