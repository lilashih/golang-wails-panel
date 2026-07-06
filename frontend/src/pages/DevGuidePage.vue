<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { FileText, RefreshCw } from "lucide-vue-next";
import { GetDocument, ListDocuments } from "@bindings/gbase/src/service/readme_reader/service.js";
import AppButton from "@/components/ui/AppButton.vue";
import EmptyState from "@/components/ui/EmptyState.vue";

interface DocumentOption {
  key: string;
  title: string;
  filename: string;
}

interface DocumentContent extends DocumentOption {
  content: string;
}

const documents = ref<DocumentOption[]>([]);
const selectedKey = ref("");
const currentDocument = ref<DocumentContent | null>(null);
const loading = ref(false);
const error = ref("");

const renderedHtml = computed(() => renderMarkdown(currentDocument.value?.content || ""));

function toErrorMessage(err: unknown) {
  return err instanceof Error ? err.message : String(err);
}

function isDocumentOption(value: unknown): value is DocumentOption {
  if (!value || typeof value !== "object") {
    return false;
  }

  const doc = value as Partial<DocumentOption>;
  return typeof doc.key === "string" && typeof doc.title === "string" && typeof doc.filename === "string";
}

function normalizeDocuments(value: unknown): DocumentOption[] {
  if (!Array.isArray(value)) {
    return [];
  }

  return value.filter(isDocumentOption);
}

function isDocumentContent(value: unknown): value is DocumentContent {
  return isDocumentOption(value) && typeof (value as Partial<DocumentContent>).content === "string";
}

async function loadDocuments() {
  loading.value = true;
  error.value = "";
  try {
    const result = await ListDocuments();
    documents.value = normalizeDocuments(result);

    if (documents.value.length === 0) {
      selectedKey.value = "";
      currentDocument.value = null;
      return;
    }

    const exists = documents.value.some((doc) => doc.key === selectedKey.value);
    selectedKey.value = exists ? selectedKey.value : documents.value[0].key;
    await loadDocument(selectedKey.value);
  } catch (err) {
    error.value = toErrorMessage(err);
  } finally {
    loading.value = false;
  }
}

async function loadDocument(key: string) {
  if (!key) {
    currentDocument.value = null;
    return;
  }

  loading.value = true;
  error.value = "";
  try {
    const result = await GetDocument(key);
    if (!isDocumentContent(result)) {
      throw new Error("文件內容格式不正確");
    }
    currentDocument.value = result;
  } catch (err) {
    currentDocument.value = null;
    error.value = toErrorMessage(err);
  } finally {
    loading.value = false;
  }
}

function escapeHtml(value: string) {
  return value.replace(/[&<>"']/g, (char) => {
    const entities: Record<string, string> = {
      "&": "&amp;",
      "<": "&lt;",
      ">": "&gt;",
      '"': "&quot;",
      "'": "&#39;",
    };
    return entities[char] || char;
  });
}

function renderInline(value: string) {
  let html = escapeHtml(value);
  html = html.replace(/`([^`]+)`/g, "<code>$1</code>");
  html = html.replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");
  html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, (_match, label: string, href: string) => {
    const safeHref = escapeHtml(href.trim());
    if (/^javascript:/i.test(safeHref)) {
      return label;
    }
    return `<a href="${safeHref}" target="_blank" rel="noreferrer">${label}</a>`;
  });
  return html;
}

function renderTable(lines: string[], start: number) {
  const header = splitTableRow(lines[start]);
  const align = splitTableRow(lines[start + 1]);
  if (header.length === 0 || align.length === 0 || !align.every((cell) => /^:?-{3,}:?$/.test(cell.trim()))) {
    return null;
  }

  const rows: string[][] = [];
  let index = start + 2;
  while (index < lines.length && /^\s*\|.*\|\s*$/.test(lines[index])) {
    rows.push(splitTableRow(lines[index]));
    index += 1;
  }

  const thead = `<thead><tr>${header.map((cell) => `<th>${renderInline(cell.trim())}</th>`).join("")}</tr></thead>`;
  const tbody = `<tbody>${rows.map((row) => `<tr>${header.map((_cell, cellIndex) => `<td>${renderInline((row[cellIndex] || "").trim())}</td>`).join("")}</tr>`).join("")}</tbody>`;

  return {
    html: `<table>${thead}${tbody}</table>`,
    next: index,
  };
}

function splitTableRow(line: string) {
  return line.trim().replace(/^\|/, "").replace(/\|$/, "").split("|");
}

function renderMarkdown(markdown: string) {
  const lines = markdown.replace(/\r\n/g, "\n").split("\n");
  const html: string[] = [];
  let index = 0;
  let inCodeBlock = false;
  let codeLines: string[] = [];
  let listItems: string[] = [];
  let paragraph: string[] = [];

  function flushParagraph() {
    if (paragraph.length > 0) {
      html.push(`<p>${renderInline(paragraph.join(" "))}</p>`);
      paragraph = [];
    }
  }

  function flushList() {
    if (listItems.length > 0) {
      html.push(`<ul>${listItems.map((item) => `<li>${renderInline(item)}</li>`).join("")}</ul>`);
      listItems = [];
    }
  }

  while (index < lines.length) {
    const line = lines[index];
    const trimmed = line.trim();

    if (trimmed.startsWith("```")) {
      if (inCodeBlock) {
        html.push(`<pre><code>${escapeHtml(codeLines.join("\n"))}</code></pre>`);
        codeLines = [];
        inCodeBlock = false;
      } else {
        flushParagraph();
        flushList();
        inCodeBlock = true;
      }
      index += 1;
      continue;
    }

    if (inCodeBlock) {
      codeLines.push(line);
      index += 1;
      continue;
    }

    if (!trimmed) {
      flushParagraph();
      flushList();
      index += 1;
      continue;
    }

    if (/^\s*\|.*\|\s*$/.test(line) && index + 1 < lines.length) {
      const table = renderTable(lines, index);
      if (table) {
        flushParagraph();
        flushList();
        html.push(table.html);
        index = table.next;
        continue;
      }
    }

    const heading = /^(#{1,6})\s+(.+)$/.exec(trimmed);
    if (heading) {
      flushParagraph();
      flushList();
      const level = heading[1].length;
      html.push(`<h${level}>${renderInline(heading[2])}</h${level}>`);
      index += 1;
      continue;
    }

    const list = /^[-*+]\s+(.+)$/.exec(trimmed);
    if (list) {
      flushParagraph();
      listItems.push(list[1]);
      index += 1;
      continue;
    }

    const quote = /^>\s*(.*)$/.exec(trimmed);
    if (quote) {
      flushParagraph();
      flushList();
      html.push(`<blockquote>${renderInline(quote[1])}</blockquote>`);
      index += 1;
      continue;
    }

    paragraph.push(trimmed);
    index += 1;
  }

  flushParagraph();
  flushList();
  if (inCodeBlock) {
    html.push(`<pre><code>${escapeHtml(codeLines.join("\n"))}</code></pre>`);
  }

  return html.join("\n");
}

onMounted(() => {
  loadDocuments().catch(console.error);
});
</script>

<template>
  <section class="flex h-full min-h-0 flex-col overflow-hidden">
    <div class="flex-none border-b border-app-border bg-app-surface px-4 py-3">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div class="min-w-0">
          <p v-if="currentDocument" class="text-sm text-app-text-muted">{{ currentDocument.filename }}</p>
        </div>
        <div class="flex min-w-0 items-center gap-2">
          <label class="flex min-w-0 items-center gap-2 text-sm text-app-text-muted">
            <FileText class="h-4 w-4 shrink-0" />
            <select
              v-model="selectedKey"
              class="h-9 min-w-48 rounded-md border border-app-border bg-app-surface px-3 text-sm text-app-text outline-none focus:border-app-primary focus:ring-2 focus:ring-app-primary/20 disabled:cursor-not-allowed disabled:opacity-60"
              :disabled="loading || documents.length === 0"
              @change="loadDocument(selectedKey)"
            >
              <option v-for="doc in documents" :key="doc.key" :value="doc.key">
                {{ doc.title }}
              </option>
            </select>
          </label>
          <AppButton size="sm" :disabled="loading" @click="loadDocuments">
            <RefreshCw class="h-4 w-4" />
            重新整理
          </AppButton>
        </div>
      </div>
      <div v-if="error" class="mt-3 rounded-md border border-app-danger/30 bg-app-danger-soft px-3 py-2 text-sm text-app-danger">
        {{ error }}
      </div>
    </div>

    <div class="min-h-0 flex-1 overflow-auto p-4">
      <EmptyState
        v-if="!loading && documents.length === 0"
        title="目前沒有可閱讀的 README 文件"
        message="請確認專案根目錄是否存在 README.md。"
      />
      <EmptyState
        v-else-if="!loading && !currentDocument"
        title="尚未載入文件內容"
        message="請從上方下拉選單選擇文件。"
      />
      <div v-else class="mx-auto max-w-5xl rounded-md border border-app-border bg-app-surface px-6 py-5 shadow-sm">
        <div v-if="loading" class="py-10 text-center text-sm text-app-text-muted">讀取中...</div>
        <article v-else class="markdown-body" v-html="renderedHtml" />
      </div>
    </div>
  </section>
</template>
