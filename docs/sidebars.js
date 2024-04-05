/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation
 The sidebars can be generated from the filesystem, or explicitly defined here.
 Create as many sidebars as you want.
 */

module.exports = {
  // By default, Docusaurus generates a sidebar from the docs folder structure
  documentationSidebar: [
    {

    },
    'intro',
    {
      type: 'category',
      label: 'Quick Start',
      link: {
        type: 'generated-index',
        slug: '/getting-started',
        title: 'Quick Start',
        description: "Get Started with Bacalhau!",
      },
      collapsed: false,
      items: [
        'getting-started/installation',
        'getting-started/architecture',
        'getting-started/docker-workload-onboarding',
        'getting-started/wasm-workload-onboarding',
        'getting-started/resources'
      ],
    },
    {
      type: 'category',
      label: 'How-to Guides',
      link: {
        type: 'generated-index',
        title: 'Setting Up',
        slug: '/setting-up',
      },
      collapsed: true,
      items: [
        {
          type: 'autogenerated',
          dirName: 'setting-up',
        },
      ],
    },
    {
      type: 'category',
      label: 'References',
      link: {
        type: 'generated-index',
        title: 'Advanced Guides',
        slug: '/advanced-guides',
      },
      collapsed: true,
      items: [
        {
          type: 'autogenerated',
          dirName: 'dev',
        },
      ]
    },
    {
      type: 'category',
      label: 'Examples',
      link: {
        type: 'generated-index',
        title: 'Examples',
        slug: '/examples',
        description: "Bacalhau comes pre-loaded with exciting examples to showcase its abilities and help get you started.",
      },
      collapsed: true,
      items: [
        {
          type: 'category',
          label: 'Data Engineering',
          link: {
            type: 'generated-index',
            description: "This directory contains examples relating to data engineering workloads. The goal is to provide a range of examples that show you how to work with Bacalhau in a variety of use cases.",
          },
          items: [
            'examples/data-engineering/blockchain-etl/index',
            'examples/data-engineering/csv-to-avro-or-parquet/index',
            'examples/data-engineering/DuckDB/index',
            'examples/data-engineering/image-processing/index',
            'examples/data-engineering/oceanography-conversion/index',
            'examples/data-engineering/simple-parallel-workloads/index',
          ],
        },
        {
          type: 'category',
          label: 'Model Inference',
          link: {
            type: 'generated-index',
            description: "This directory contains examples relating to model inference workloads.",
          },
          items: [
            'examples/model-inference/Huggingface-Model-Inference/index',
            'examples/model-inference/object-detection-yolo5/index',
            'examples/model-inference/S3-Model-Inference/index',
            'examples/model-inference/Stable-Diffusion-CKPT-Inference/index',
            'examples/model-inference/stable-diffusion-cpu/index',
            'examples/model-inference/stable-diffusion-gpu/index',
            'examples/model-inference/StyleGAN3/index',
            'examples/model-inference/EasyOCR/index',
            'examples/model-inference/Openai-Whisper/index',
          ],
        },
        {
          type: 'category',
          label: 'Model Training',
          link: {
            type: 'generated-index',
            description: "This directory contains examples relating to model training workloads.",
          },
          items: [
            'examples/model-training/Stable-Diffusion-Dreambooth/index',
            'examples/model-training/Training-Pytorch-Model/index',
            'examples/model-training/Training-Tensorflow-Model/index',
          ],
        },
        {
          type: 'category',
          label: 'Molecular Dynamics',
          link: {
            type: 'generated-index',
            description: "This directory contains examples relating to performing common tasks with Bacalhau.",
          },
          items: [
            'examples/molecular-dynamics/BIDS/index',
            'examples/molecular-dynamics/Coreset/index',
            'examples/molecular-dynamics/Genomics/index',
            'examples/molecular-dynamics/Gromacs/index',
            'examples/molecular-dynamics/openmm/index',
          ],
        },
      ],
    },
    {
      type: 'category',
      label: 'Integrations',
      link: {
        type: 'generated-index',
        title: 'Integrations',
        slug: '/integration',
      },
      collapsed: true,
      items: [
        'integration/python-sdk',
        'integration/apache-airflow',
        'integration/lilypad',
        'integration/wasm-observe'
      ],
    },
    'faqs',
    {
      type: 'category',
      label: 'Community',
      link: {
        type: 'generated-index',
        title: 'Community',
        slug: '/community',
      },
      collapsed: true,
      items: [
        'community/compute-landscape',
        'community/style-guide',
        'community/ways-to-contribute',
      ],
    }
  ]
}
