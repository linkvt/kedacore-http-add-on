package util

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Signaler", func() {
	Context("New", func() {
		It("returns a signaler of capacity 1", func() {
			i := NewSignaler()

			s, ok := i.(signaler)
			Expect(ok).To(BeTrue())

			c := cap(s)
			Expect(c).To(Equal(1))
		})
	})

	Context("Signal", func() {
		It("produces on channel", func() {
			s := make(signaler, 1)

			s.Signal()

			select {
			case <-s:
			default:
				Fail("channel should not be empty")
			}
		})

		It("does not block when channel is full", func() {
			s := make(signaler, 1)
			s <- struct{}{}

			s.Signal()
		})
	})

	Context("Wait", func() {
		It("returns nil when channel is not empty", func() {
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
			defer cancel()

			s := make(signaler, 1)
			s <- struct{}{}

			err := s.Wait(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns err when context is done", func() {
			ctx, cancel := context.WithCancel(context.TODO())
			cancel()

			s := make(signaler, 1)

			err := s.Wait(ctx)
			Expect(err).To(MatchError(context.Canceled))
		})
	})

	Context("E2E", func() {
		It("succeeds", func() {
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()

			s := NewSignaler()

			s.Signal()
			s.Signal()

			timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second)
			defer timeoutCancel()

			err := s.Wait(timeoutCtx)
			Expect(err).NotTo(HaveOccurred())

			cancel()

			err = s.Wait(ctx)
			Expect(err).To(MatchError(context.Canceled))
		})
	})
})
