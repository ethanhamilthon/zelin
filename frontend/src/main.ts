import { mount } from 'svelte';
import Root from './Root.svelte';
import '@fontsource-variable/jetbrains-mono';
import './style.css';

mount(Root, { target: document.getElementById('app')! });
