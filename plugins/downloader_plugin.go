package plugins

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func init() {
	Register(Command{
		Name:  "tiktok",
		Tags:  []string{"downloader"},
		Help:  "Mengunduh video TikTok tanpa tanda air (watermark)",
		Limit: true,
		Execute: runTiktokDownloader,
	})
}

func runTiktokDownloader(ctx *Context) error {
	if ctx.Query == "" {
		return ctx.Reply("⚠️ Harap sertakan URL video TikTok! Contoh: `.tiktok https://vt.tiktok.com/...`")
	}

	if !strings.Contains(ctx.Query, "tiktok.com") {
		return ctx.Reply("❌ Harap berikan link TikTok yang valid!")
	}

	ctx.React("🕛")
	ctx.Reply("⏳ Sedang memproses dan mengunduh video, mohon tunggu...")

	// 1. Dapatkan video bytes dari API Downloader.
	// Ini adalah template terintegrasi. Sebagai contoh, kita men-download video contoh secara langsung.
	// Pada implementasi riil, Anda bisa melakukan request API downloader pilihan Anda untuk mendapatkan video link yang bersih.
	// Contoh: resp, err := http.Get("https://api.tiktok-downloader/v1?url=" + ctx.Query)
	
	videoUrl := "https://files.catbox.moe/opmcq5.mp4" // Video dummy sebagai demo media sending
	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Get(videoUrl)
	if err != nil {
		return fmt.Errorf("gagal mengunduh media dari server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status response server tidak OK: %d", resp.StatusCode)
	}

	mediaData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("gagal membaca byte media: %v", err)
	}

	// 2. Upload byte media ke server WhatsApp menggunakan whatsmeow Client.
	uploadResp, err := ctx.Client.Upload(context.Background(), mediaData, whatsmeow.MediaVideo)
	if err != nil {
		return fmt.Errorf("gagal mengunggah media ke server WhatsApp: %v", err)
	}

	// 3. Bangun WhatsApp Video Message.
	videoMsg := &waE2E.VideoMessage{
		URL:           proto.String(uploadResp.URL),
		DirectPath:    proto.String(uploadResp.DirectPath),
		MediaKey:      uploadResp.MediaKey,
		Mimetype:      proto.String("video/mp4"),
		FileLength:    proto.Uint64(uint64(len(mediaData))),
		FileSHA256:    uploadResp.FileSHA256,
		FileEncSHA256: uploadResp.FileEncSHA256,
		Caption:       proto.String(fmt.Sprintf("🎥 *TikTok Downloader*\n\nRequest oleh: %s", ctx.PushName)),
	}

	msg := &waE2E.Message{
		VideoMessage: videoMsg,
	}

	// 4. Kirim pesan video ke chat.
	_, err = ctx.Client.SendMessage(context.Background(), ctx.Event.Info.Chat, msg)
	if err != nil {
		return fmt.Errorf("gagal mengirimkan video ke chat: %v", err)
	}

	ctx.React("✅")
	return nil
}
