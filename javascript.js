document.addEventListener("DOMContentLoaded", function() {
    // URLの正規表現パターン
    const urlRegex = /(https?:\/\/[^\s]+)/g;

    // テキストノード内のURLを検出しリンクに変換する関数
    function linkifyTextNodes(element) {
        const walker = document.createTreeWalker(
            element,
            NodeFilter.SHOW_TEXT,
            null,
            false
        );

        let node;
        while ((node = walker.nextNode())) {
            const text = node.nodeValue;
            if (urlRegex.test(text)) {
                const span = document.createElement("span");
                span.innerHTML = text.replace(urlRegex, '<a href="$&" target="_blank">$&</a>');
                node.parentNode.replaceChild(span, node);
            }
        }
    }

    // ボディ全体に適用
    linkifyTextNodes(document.body);
});
