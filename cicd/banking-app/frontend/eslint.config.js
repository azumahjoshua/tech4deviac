export default [
    {
        // Environment settings
        languageOptions: {
            ecmaVersion: 2021,
            sourceType: 'module',
            globals: {
                browser: true,
                node: true,
            },
        },

        // Rules
        rules: {
            'no-unused-vars': 'warn',
            // Add more custom rules here:
            // 'indent': ['error', 2],
            // 'quotes': ['error', 'single'],
        },
    },
];