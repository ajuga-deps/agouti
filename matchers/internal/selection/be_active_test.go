package selection_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/matchers/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/selection"
)

var _ = Describe("BeActiveMatcher", func() {
	var (
		matcher   *BeActiveMatcher
		selection *mocks.Selection
	)

	BeforeEach(func() {
		selection = &mocks.Selection{}
		selection.StringCall.ReturnString = "CSS: #selector"
		matcher = &BeActiveMatcher{}
	})

	Describe("#Match", func() {
		Context("when the actual object is a selection", func() {
			Context("when the element is active", func() {
				It("should successfully return true", func() {
					selection.ActiveCall.ReturnActive = true
					success, err := matcher.Match(selection)
					Expect(success).To(BeTrue())
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when the element is not active", func() {
				It("should successfully return false", func() {
					selection.ActiveCall.ReturnActive = false
					success, err := matcher.Match(selection)
					Expect(success).To(BeFalse())
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when determining whether the element is active fails", func() {
				It("should return an error", func() {
					selection.ActiveCall.Err = errors.New("some error")
					_, err := matcher.Match(selection)
					Expect(err).To(MatchError("some error"))
				})
			})
		})

		Context("when the actual object is not a selection", func() {
			It("should return an error", func() {
				_, err := matcher.Match("not a selection")
				Expect(err).To(MatchError("BeActive matcher requires a Selection.  Got:\n    <string>: not a selection"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message", func() {
			selection.ActiveCall.ReturnActive = false
			matcher.Match(selection)
			message := matcher.FailureMessage(selection)
			Expect(message).To(Equal("Expected selection 'CSS: #selector' to be active"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message", func() {
			selection.ActiveCall.ReturnActive = true
			matcher.Match(selection)
			message := matcher.NegatedFailureMessage(selection)
			Expect(message).To(Equal("Expected selection 'CSS: #selector' not to be active"))
		})
	})
})
