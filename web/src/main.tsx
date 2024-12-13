import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import {createBrowserRouter, RouterProvider} from "react-router-dom";
import {ConfirmationPage} from "./ConfirmationPage.tsx";

const router = createBrowserRouter([
	{
		path: '/',
		element: <App />
	},
	{
		path: '/confirm/:token',
		element: <ConfirmationPage />
	},
])

createRoot(document.getElementById('root')!).render(
	<StrictMode>
		<RouterProvider router={router} />
	</StrictMode>,
)
