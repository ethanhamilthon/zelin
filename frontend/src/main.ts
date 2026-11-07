import { mount } from 'svelte';
import App from './App.svelte';
import '@fontsource-variable/jetbrains-mono';
import './style.css';

mount(App, { target: document.getElementById('app')! });
