/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "goqueue",
  tagline: "Redis-backed Go SDK for background tasks",
  favicon: "img/favicon.ico",
  url: "https://example.com",
  baseUrl: "/",
  onBrokenLinks: "warn",
  onBrokenMarkdownLinks: "warn",
  trailingSlash: false,
  organizationName: "newton-school",
  projectName: "goqueue",

  presets: [
    [
      "classic",
      {
        docs: {
          path: "docs",
          routeBasePath: "/",
          sidebarPath: require.resolve("./sidebars.js"),
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      },
    ],
  ],

  themeConfig: {
    navbar: {
      title: "goqueue",
      items: [
        {
          type: "doc",
          docId: "intro",
          position: "left",
          label: "Documentation",
        },
        {
          href: "https://github.com/Newton-School/goqueue",
          label: "GitHub",
          position: "right",
        },
      ],
    },
    footer: {
      style: "dark",
      copyright: `Copyright © ${new Date().getFullYear()} goqueue.`,
    },
  },
};

module.exports = config;
