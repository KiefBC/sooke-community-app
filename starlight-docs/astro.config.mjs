// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	integrations: [
		starlight({
			title: 'Sooke Community App',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/KiefBC/sooke_app' }],
			sidebar: [
				{
					label: 'Project',
					items: [
						{ label: 'Project Plan', slug: 'project-plan' },
						{ label: 'Database Schema', slug: 'database-schema' },
						{ label: 'Style Guide', slug: 'style-guide' },
						{ label: 'Common Commands', slug: 'common-commands' },
						{ label: 'Future Ideas and Alternatives', slug: 'future-ideas-and-alternatives' },
						{ label: 'Planning Discussion', slug: 'planning-discussion' },
					],
				},
				{
					label: 'Guides',
					items: [
						{ label: 'SwiftUI Data Flow', slug: 'swiftui-data-guide' },
					],
				},
				{
					label: 'Architecture Decision Records',
					autogenerate: { directory: 'decisions' },
				},
			],
		}),
	],
});
