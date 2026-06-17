/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  goqueueSidebar: [
    "intro",
    {
      type: "category",
      label: "Getting Started",
      items: [
        "getting-started/installation",
        "getting-started/quick-start",
      ],
    },
    {
      type: "category",
      label: "Core Concepts",
      items: [
        "concepts/configuration",
        "concepts/task-model",
        "concepts/producer",
        "concepts/worker",
        "concepts/scheduler",
        "concepts/workflow",
      ],
    },
    {
      type: "category",
      label: "Operate",
      items: [
        "concepts/inspect-and-admin",
        "concepts/redis-backend",
      ],
    },
    {
      type: "category",
      label: "Reference",
      items: [
        "reference/cli",
        "reference/errors",
      ],
    },
  ],
};

module.exports = sidebars;
