package parser

import (
	"errors"
	"strings"
	"testing"
)

func TestParseItemDescriptionPage(t *testing.T) {
	t.Parallel()

	raw := `<!doctype html><html><body><article>
		<section class="main-info-column">
			<div class="js-market-item-detail-description description">
				<p class="autolink whitespace-pre-line">冒頭本文</p>
			</div>
		</section>
	</article></body></html>`

	description, err := ParseItemDescriptionPage(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("ParseItemDescriptionPage() error = %v", err)
	}

	if !strings.Contains(description, "冒頭本文") {
		t.Fatalf("description = %q", description)
	}
}

func TestParseItemDescriptionPageIncludesSupplementalSections(t *testing.T) {
	t.Parallel()

	raw := `<!doctype html><html><body><article>
		<section class="main-info-column">
			<div class="js-market-item-detail-description description">
				<p class="autolink whitespace-pre-line">冒頭本文</p>
			</div>
		</section>
		<section class="shop__text">
			<h2>見出し</h2>
			<p class="js-autolink whitespace-pre-line">後続本文</p>
			<p class="js-autolink whitespace-pre-line">補足段落</p>
		</section>
	</article></body></html>`

	description, err := ParseItemDescriptionPage(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("ParseItemDescriptionPage() error = %v", err)
	}

	if !strings.Contains(description, "冒頭本文") {
		t.Fatalf("description = %q", description)
	}
	if !strings.Contains(description, "見出し\n後続本文") {
		t.Fatalf("description = %q", description)
	}
	if !strings.Contains(description, "補足段落") {
		t.Fatalf("description = %q", description)
	}
}

func TestParseItemDescriptionPageNotFound(t *testing.T) {
	t.Parallel()

	_, err := ParseItemDescriptionPage(strings.NewReader(`<!doctype html><html><body></body></html>`))
	if err == nil {
		t.Fatal("ParseItemDescriptionPage() error = nil, want error")
	}
	if !errors.Is(err, err) {
		t.Fatalf("ParseItemDescriptionPage() error = %v", err)
	}
}

func TestParseItemAvatarsPage(t *testing.T) {
	t.Parallel()

	raw := `<!doctype html><html><body>
		<ul id="variations">
			<li><div class="variation-name">フルパック - FullPack</div></li>
			<li><div class="variation-name">ルミナ - LUMINA</div></li>
			<li><div class="variation-name">ショコラ - Chocolat ※共通素体あり</div></li>
		</ul>
	</body></html>`

	avatars, err := ParseItemAvatarsPage(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("ParseItemAvatarsPage() error = %v", err)
	}

	if len(avatars) != 3 {
		t.Fatalf("len(avatars) = %d, want 3", len(avatars))
	}
	if avatars[0] != "フルパック - FullPack" || avatars[1] != "ルミナ - LUMINA" || avatars[2] != "ショコラ - Chocolat ※共通素体あり" {
		t.Fatalf("avatars = %#v", avatars)
	}
}
