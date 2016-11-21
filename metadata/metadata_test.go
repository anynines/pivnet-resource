package metadata_test

import (
	"fmt"

	"github.com/anynines/pivnet-resource/metadata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metadata", func() {
	Describe("Validate", func() {
		var data metadata.Metadata
		BeforeEach(func() {
			data = metadata.Metadata{
				Release: &metadata.Release{
					Version:     "1.0.0",
					ReleaseType: "All In One",
					EULASlug:    "some-other-eula",
				},
				ProductFiles: []metadata.ProductFile{
					{File: "hello.txt", Description: "available"},
				},
			}
		})

		Context("when release is missing", func() {
			BeforeEach(func() {
				data.Release = nil
			})

			It("returns an error", func() {
				Expect(data.Validate()).To(MatchError(fmt.Sprintf("missing required value %q", "release")))
			})
		})

		Context("when eula slug is missing", func() {
			BeforeEach(func() {
				data.Release.EULASlug = ""
			})

			It("returns an error", func() {
				Expect(data.Validate()).To(MatchError(fmt.Sprintf("missing required value %q", "eula_slug")))
			})
		})

		Context("when version is missing", func() {
			BeforeEach(func() {
				data.Release.Version = ""
			})

			It("returns an error", func() {
				Expect(data.Validate()).To(MatchError(fmt.Sprintf("missing required value %q", "version")))
			})
		})

		Context("when release type is missing", func() {
			BeforeEach(func() {
				data.Release.ReleaseType = ""
			})

			It("returns an error", func() {
				Expect(data.Validate()).To(MatchError(fmt.Sprintf("missing required value %q", "release_type")))
			})
		})

		Context("when product files are missing", func() {
			BeforeEach(func() {
				data.ProductFiles[0].File = ""
			})

			It("returns an error", func() {
				Expect(data.Validate()).To(MatchError("empty value for file"))
			})
		})

		Context("when dependencies exist with id 0", func() {
			BeforeEach(func() {
				data.Dependencies = []metadata.Dependency{
					{
						Release: metadata.DependentRelease{
							ID:      0,
							Version: "abcd",
							Product: metadata.Product{
								Slug: "some-product",
							},
						},
					},
				}
			})

			It("returns without error", func() {
				Expect(data.Validate()).NotTo(HaveOccurred())
			})

			Context("when release version is empty", func() {
				BeforeEach(func() {
					data.Dependencies[0].Release.Version = ""
				})

				It("returns an error", func() {
					err := data.Validate()
					Expect(err).To(HaveOccurred())

					Expect(err.Error()).To(MatchRegexp(".*dependency\\[0\\]"))
				})
			})

			Context("when product slug is empty", func() {
				BeforeEach(func() {
					data.Dependencies[0].Release.Product.Slug = ""
				})

				It("returns an error", func() {
					err := data.Validate()
					Expect(err).To(HaveOccurred())

					Expect(err.Error()).To(MatchRegexp(".*dependency\\[0\\]"))
				})
			})
		})

		Context("when upgrade paths are provided", func() {
			BeforeEach(func() {
				data.UpgradePaths = []metadata.UpgradePath{
					{
						ID:      1234,
						Version: "abcd",
					},
				}
			})

			It("returns without error", func() {
				Expect(data.Validate()).NotTo(HaveOccurred())
			})

			Context("when id is non-zero and version is empty", func() {
				BeforeEach(func() {
					data.UpgradePaths[0].Version = ""
				})

				It("returns without error", func() {
					Expect(data.Validate()).NotTo(HaveOccurred())
				})
			})

			Context("when id is 0 and version is non-empty", func() {
				BeforeEach(func() {
					data.UpgradePaths[0].ID = 0
				})

				It("returns without error", func() {
					Expect(data.Validate()).NotTo(HaveOccurred())
				})
			})

			Context("when id is 0 and version is empty", func() {
				BeforeEach(func() {
					data.UpgradePaths[0].ID = 0
					data.UpgradePaths[0].Version = ""
				})

				It("returns an error", func() {
					err := data.Validate()
					Expect(err).To(HaveOccurred())

					Expect(err.Error()).To(MatchRegexp(".*upgrade_paths\\[0\\]"))
				})
			})
		})
	})
})
